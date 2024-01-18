package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sdecay/podcasty/internal/database"
)

func scrape(db *database.Queries, threads int, delay time.Duration) {
	ticker := time.NewTicker(delay)

	log.Printf("Scraping on %v threads every %d seconds...\n", threads, delay/time.Second)

	// empty for so it starts immediately on program run
	for ; ; <-ticker.C {
		existing := getExistingPosts(db)

		feeds, err := db.FetchNextFeeds(context.Background(), int32(threads))
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wtGrp := &sync.WaitGroup{}

		for _, feed := range feeds {
			wtGrp.Add(1)
			go scrapeFeed(db, wtGrp, feed, existing)
		}
		wtGrp.Wait()
	}
}

func scrapeFeed(db *database.Queries, wtGrp *sync.WaitGroup, feed database.Feed, existing map[string]bool) {
	fmt.Println("scraping feed:", feed.Name)
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
		if existing[item.Link] {
			log.Println("post exists in database, skipping:", item.Link)
			continue
		}

		description, pubTime := normalizePostData(item.Description, item.PubDate)

		fmt.Println(description, pubTime)

		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubTime,
				Url:         item.Link,
				FeedID:      feed.ID,
			})
		if err != nil {
			log.Println("failed to create post:", err)
		}
	}
	log.Printf("feed %s scraped with %d posts", feed.Name, len(rssFeed.Channel.Item))
}

func getExistingPosts(db *database.Queries) map[string]bool {
	// grab existing urls (the only unique, searchable thing in table)
	// probably not best to store this in memory at scale but hey
	urls, err := db.GetPostUniques(context.Background())
	if err != nil {
		fmt.Println("problem fetching uniques uh oh!")
	}
	context.Background().Done() // this was just a dumb guess.  worked!

	existing := map[string]bool{}

	for _, url := range urls {
		existing[url] = true
	}
	return existing
}

func normalizePostData(description string, pubDate string) (sql.NullString, time.Time) {
	newDescr := sql.NullString{}

	if description != "" {
		newDescr.String = description
		newDescr.Valid = true
	}

	newTime, err := time.Parse(time.RFC1123Z, pubDate)
	if err != nil {
		log.Printf("could not parse funky blog item date (%v): %s", pubDate, err)
	}

	return newDescr, newTime
}
