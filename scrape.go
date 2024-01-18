package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sdecay/podcasty/internal/database"
)

func scrape(db *database.Queries, threads int, delay time.Duration) {
	ticker := time.NewTicker(delay)

	log.Printf("Scraping on %v threads every %d...\n", threads, delay)

	// empty for so it starts immediately on program run
	for ; ; <-ticker.C {
		fmt.Println("in ticker")
		feeds, err := db.FetchNextFeeds(context.Background(), int32(threads))
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wtGrp := &sync.WaitGroup{}

		for _, feed := range feeds {
			wtGrp.Add(1)

			go scrapeFeed(db, wtGrp, feed)
		}
		wtGrp.Wait()
	}
}

func scrapeFeed(db *database.Queries, wtGrp *sync.WaitGroup, feed database.Feed) {
	defer wtGrp.Done()

	_, err := db.MarkFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking fetched", err)
	}

	rssFeed, err := UrlToRssFeed(feed.Url)
	if err != nil {
		log.Printf("could not parse rss feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("found post:", item.Title)
	}

	log.Printf("feed %s scraped with %d posts", feed.Name, len(rssFeed.Channel.Item))
}
