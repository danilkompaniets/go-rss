package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// Определяем структуру для парсинга RSS
type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		Items       []Item `xml:"item"` // Исправлено: теперь это массив объектов
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RssFeed, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RssFeed{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RssFeed{}, err
	}

	var rssFeed RssFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RssFeed{}, err
	}

	return rssFeed, nil
}
