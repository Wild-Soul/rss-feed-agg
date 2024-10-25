package worker

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
	"github.com/Wild-Soul/go-rss-feed-agg/rss"
	"github.com/google/uuid"
)

func StartScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C { // start executing immediately.

		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Println("Error: fetching feeds:", err.Error())
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, feed, wg)
		}
		wg.Wait()

		log.Printf("Processed %v feeds\n", concurrency)
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err.Error())
		return
	}

	rssFeed, err := rss.UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err.Error())
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		publishDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Could't parse date: %v, err: %v\n", item.PubDate, err.Error())
			continue
		}

		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: publishDate,
				Url:         item.Link,
				FeedID:      feed.ID,
			},
		)

		if err != nil {
			log.Printf("Failed to save post: %v from feed: %v\n", item.Title, feed.ID)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.ID, len(rssFeed.Channel.Item))
}
