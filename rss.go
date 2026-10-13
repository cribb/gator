package main

import (
	// "bufio"
	// "encoding/json"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	// "strings"
	//"time"
	// "github.com/cribb/gator/internal/config"
	// "github.com/cribb/gator/internal/database"
	// "github.com/google/uuid"
	// "go.opentelemetry.io/otel/semconv/v1.43.0/dbconv"
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

// 	// res :=
// 	xml.Unmarshal(data, &res)
// 	// ctx.Value()
// 	foo := xml.Comment("example comment")

// 	bar := html.EscapeString("")
// 	bar = strings.ToLower(bar)
// 	fmt.Printf("%v: %v\n", foo, bar)

// 	// return nil, fmt.Errorf("Not yet implemented.")
// }

// func getDataFromNetwork(url string) ([]byte, error) {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		fmt.Printf("Error in GET functionality.\n")
// 		return nil, fmt.Errorf("error creating request: %w", err)
// 	}

// 	body, err := io.ReadAll(res.Body)
// 	// fmt.Printf(" --> body: %v, err: %v\n", res.Body, err)
// 	defer res.Body.Close()

// 	if err != nil {
// 		return nil, err
// 	}

// 	if res.StatusCode > 299 {
// 		// fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
// 		return nil, fmt.Errorf("Response failed with status code: %d", res.StatusCode)
// 	}

// 	return body, nil
// }

// // Pulls data either directly from the website API (slow) or from the local cache (fast)
// func fetchData(config *cliConfig, url string) ([]byte, error) {
// 	var body []byte
// 	var ok bool
// 	var err error

// 	body, ok = config.cache.Get(url)
// 	if !ok {
// 		body, err = getDataFromNetwork(url)

// 		if err != nil {
// 			return nil, err
// 		}
// 		config.cache.Add(url, body)
// 	}
// 	return body, nil

// }
