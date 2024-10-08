package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/jsec/gator/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*Feed, error) {
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

func ScrapeFeeds(ctx context.Context, db *database.Queries, timeBetweenReqs string) {
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		fmt.Println("Invalid duration string:", err.Error())
		return
	}

	fmt.Println("Scraping feeds every", timeBetweenReqs)

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		feed, err := db.GetNextFeedToFetch(ctx)
		if err != nil {
			fmt.Println("Error retrieving feed to fetch:", err.Error())
			return
		}

		fetched, err := fetchFeed(ctx, feed.Url)
		if err != nil {
			fmt.Println("Error fetching feed:", err.Error())
			return
		}

		for _, item := range fetched.Channel.Item {
			fmt.Println("- ", item.Title)
		}
	}
}
