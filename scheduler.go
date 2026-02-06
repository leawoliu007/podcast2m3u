package main

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/nsoufr/podfeed"
)

func StartScheduler(config Config) {
	c := cron.New()

	for _, sub := range config.Subscriptions {
		sub := sub // capture variable for closure
		
		schedule := config.Global.UpdateInterval
		if sub.Cron != "" {
			schedule = sub.Cron
		}

		if schedule == "" {
			log.Printf("No schedule found for %s, skipping auto-update", sub.Name)
			continue
		}

		_, err := c.AddFunc(schedule, func() {
			log.Printf("Updating subscription: %s", sub.Name)
			processSubscription(sub)
		})

		if err != nil {
			log.Printf("Error adding cron job for %s: %v", sub.Name, err)
		} else {
			log.Printf("Scheduled %s with %s", sub.Name, schedule)
		}
	}

	c.Start()
	log.Println("Scheduler started...")
	select {} // Block forever
}

func processSubscription(sub Subscription) {
	podcast, err := podfeed.Fetch(context.Background(), sub.URL)
	if err != nil {
		log.Printf("Failed to fetch podcast %s: %v", sub.Name, err)
		return
	}

	// Update DB (Simple insert/update for now)
	dbPodcast := Podcast{
		Name: sub.Name,
		URL:  sub.URL,
		LastUpdated: podcast.PubDate, // Assuming PubDate matches format or we just store string
	}
	
	// FirstOrCreate fits well here to avoid duplicates
	result := DB.Where(Podcast{URL: sub.URL}).FirstOrCreate(&dbPodcast)
	if result.Error != nil {
		log.Printf("Failed to update DB for %s: %v", sub.Name, result.Error)
	}

	// Update M3U
	err = M3u(podcast, sub.OutputPath)
	if err != nil {
		log.Printf("Failed to write M3U for %s: %v", sub.Name, err)
	} else {
		log.Printf("Successfully updated %s", sub.Name)
	}
}
