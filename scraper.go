package main

import (
	"context"
	"database/sql"
	"github.com/danilkompaniets/go-rss/internal/database"
	"github.com/google/uuid"
	"log"
	"strings"
	"sync"
	"time"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Starting scraping with %v connections every %v seconds", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching next feed from database:", err)
			continue
		}

		wg := sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, &wg, feed)
		}

		wg.Wait()

	}

}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Println("error marking feed as fetched:", err)
	}

	rssFeed, err := urlToFeed(feed.Url)

	if err != nil {
		log.Println("error parsing feed url:", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		desc := sql.NullString{}

		if item.Description != "" {
			desc = sql.NullString{String: item.Description, Valid: true}
		}
		t, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Println("error parsing pub date:", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Name:        item.Title,
			Description: desc,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Println("error creating post:", err)
		}
	}

}
