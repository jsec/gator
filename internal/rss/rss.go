package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func GetFeed(ctx context.Context, feedURL string) (*Feed, error) {
	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error generating HTTP request:", err)
	}

	req.Header.Add("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error fetching RSS feed:", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body:", err)
	}

	var feed Feed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling XML feed", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
