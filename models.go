package main

import (
	"github.com/danilkompaniets/go-rss/internal/database"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(u database.User) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		ApiKey:    u.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(f database.Feed) Feed {
	return Feed{
		ID:        f.ID,
		UpdatedAt: f.UpdatedAt,
		CreatedAt: f.CreatedAt,
		Name:      f.Name,
		URL:       f.Url,
		UserID:    f.UserID,
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	var resultFeeds []Feed

	for _, dbFeed := range feeds {
		resultFeeds = append(
			resultFeeds,
			databaseFeedToFeed(dbFeed),
		)
	}

	return resultFeeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowToFeedFollow(f database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        f.ID,
		UpdatedAt: f.UpdatedAt,
		CreatedAt: f.CreatedAt,
		UserID:    f.UserID,
		FeedID:    f.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	var resultFeedFollows []FeedFollow

	for _, dbFeedFollow := range dbFeedFollows {
		resultFeedFollows = append(resultFeedFollows, databaseFeedFollowToFeedFollow(dbFeedFollow))
	}

	return resultFeedFollows
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(d database.Post) Post {
	var desc *string

	if d.Description.Valid {
		desc = &d.Description.String
	}

	return Post{
		ID:          d.ID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		Name:        d.Name,
		Description: desc,
		PublishedAt: d.PublishedAt,
		Url:         d.Url,
		FeedID:      d.FeedID,
	}
}

func databasePostsToPosts(posts []database.Post) []Post {
	var resultPosts []Post

	for _, dbPost := range posts {
		resultPosts = append(resultPosts, databasePostToPost(dbPost))
	}

	return resultPosts
}
