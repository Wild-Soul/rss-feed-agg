package worker

import (
	"context"
	"database/sql"
	"fmt"
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
	log.Printf("Scraping with %v goroutines every %s", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			feeds, err := db.GetNextFeedsToFetch(ctx, database.GetNextFeedsToFetchParams{
				LastFetchedAt: time.Now().UTC(),
				Limit:         int32(concurrency),
			})

			if err != nil {
				log.Println("Error fetching feeds:", err)
				continue
			}

			var wg sync.WaitGroup
			sem := make(chan struct{}, concurrency) // Semaphore, using channels, for concurrency control.

			for _, feed := range feeds {
				sem <- struct{}{} // Acquire a token
				wg.Add(1)
				go func(feed database.Feed) {
					defer wg.Done()
					defer func() { <-sem }() // Release the token

					if err := scrapeFeed(ctx, db, feed); err != nil {
						log.Println("Error scraping feed:", err)
					}
				}(feed)
			}
			wg.Wait()
			log.Printf("Processed %v feeds\n", len(feeds))
		}
	}
}

func scrapeFeed(ctx context.Context, db *database.Queries, feed database.Feed) error {
	_, err := db.MarkFeedAsFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("marking feed as fetched: %w", err)
	}

	rssFeed, err := rss.UrlToFeed(feed.Url)
	if err != nil {
		return fmt.Errorf("fetching feed: %w", err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse date: %v, error: %v\n", item.PubDate, err)
			continue
		}

		if _, err := db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: publishDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		}); err != nil {
			log.Printf("Failed to save post: %v from feed: %v, error: %v\n", item.Title, feed.ID, err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.ID, len(rssFeed.Channel.Item))
	return nil
}
