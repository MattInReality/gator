package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
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
