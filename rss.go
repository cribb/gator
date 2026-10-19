package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/cribb/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	var feed RSSFeed

	// initialize the request and set basic params
	// http.MethodGet = "GET"
	// fourth parameter "body" (data sent with request). Is nil for GET request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Bad request.\n  (%w)\n", err)
	}

	// customize params on the request (req)
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Bad response.\n  (%w)\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error in reading response.\n  (%w)\n", err)
	}

	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	// Tidy up the legit strings so the user isn't made to read escaped html (yuck)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {

	// func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error
	// func (q *Queries) GetNextFeedToFetch(ctx context.Context) (Feed, error)

	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	// fmt.Println("Item Titles:")
	// for _, item := range rssFeed.Channel.Item {
	// 	fmt.Printf(" * %v\n", item.Title)
	// }

	fmt.Println("Item Titles:")
	for _, item := range rssFeed.Channel.Item {

		now := time.Now().UTC()

		var sqlDesc sql.NullString
		if len(item.Description) != 0 {
			sqlDesc.String = item.Description
			sqlDesc.Valid = true
		}

		var pubTime sql.NullTime
		t, err := time.Parse(time.RFC1123, item.PubDate)
		if err == nil {
			pubTime.Time = t
			pubTime.Valid = true
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: sqlDesc,
			PublishedAt: pubTime,
			FeedID:      feed.ID,
		}

		s.db.CreatePost(context.Background(), params)
	}

	return nil
}
