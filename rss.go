package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/MattInReality/gator/internal/database"
	"github.com/google/uuid"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	feed := &RSSFeed{}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return feed, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return feed, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return feed, err
	}
	err = xml.Unmarshal(data, feed)
	if err != nil {
		return feed, err
	}
	feed = unescapeData(feed)
	return feed, nil
}

func unescapeData(data *RSSFeed) *RSSFeed {
	data.Channel.Description = html.UnescapeString(data.Channel.Description)
	data.Channel.Title = html.UnescapeString(data.Channel.Title)
	for i := 0; i < len(data.Channel.Item); i++ {
		data.Channel.Item[i].Description = html.UnescapeString(data.Channel.Item[i].Description)
		data.Channel.Item[i].Title = html.UnescapeString(data.Channel.Item[i].Title)
	}
	return data
}

func scrapeFeeds(s *state) error {
	f, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), f.Url)
	if err != nil {
		return err
	}
	// NullTime requires the Valid true is set. If valid false is set, SQL set's null.
	feedUpdate := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
		ID:            f.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), feedUpdate)
	if err != nil {
		return err
	}
	for _, item := range feed.Channel.Item {
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		valid := true
		if err != nil {
			fmt.Printf("%v\n", err)
			valid = false
		}
		np := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubAt, Valid: valid},
			FeedID:      f.ID,
		}
		_, err = s.db.CreatePost(context.Background(), np)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				fmt.Printf("error creating post: %v\n", err)
			}
		}
	}
	return nil
}
