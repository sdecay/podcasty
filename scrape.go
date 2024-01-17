package main

import (
	"log"
	"time"

	"github.com/sdecay/podcasty/internal/database"
)

func scrape(db *database.Queries, threads int, delay time.Duration) {
	log.Printf("Scraping on %v threads every %d...\n", threads, delay)

	ticker := time.NewTicker(delay)

	// empty for {} so it starts immediately on program run
	for ; ; <-ticker.C {

	}
}
