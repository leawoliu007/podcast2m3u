package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"crypto/tls"
	"path/filepath"

	"github.com/robfig/cron/v3"
	"github.com/mmcdole/gofeed"
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
			processSubscription(sub, config.Global)
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

func processSubscription(sub Subscription, globalConfig GlobalConfig) {
	var feed *gofeed.Feed
	var err error
	
	fp := gofeed.NewParser()

	if globalConfig.SkipCertVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		fp.Client = &http.Client{Transport: tr}
	}
	
	feed, err = fp.ParseURL(sub.URL)

	if err != nil {
		log.Printf("Failed to fetch/parse podcast %s: %v", sub.Name, err)
		return
	}

	// Update DB (Simple insert/update for now)
	dbPodcast := Podcast{
		Name: sub.Name,
		URL:  sub.URL,
		LastUpdated: time.Now().Format(time.RFC3339),
	}
	
	// FirstOrCreate fits well here to avoid duplicates
	result := DB.Where(Podcast{URL: sub.URL}).FirstOrCreate(&dbPodcast)
	if result.Error != nil {
		log.Printf("Failed to update DB for %s: %v", sub.Name, result.Error)
	}

	// Construct Output Path: GlobalPath / Name.m3u
	// Sanitize name for filename? For now assume Name is safe or user handles it.
	filename := fmt.Sprintf("%s.m3u", sub.Name)
	outputPath := filepath.Join(globalConfig.OutputPath, filename)

	// Update M3U
	err = M3u(feed, outputPath)
	if err != nil {
		log.Printf("Failed to write M3U for %s: %v", sub.Name, err)
	} else {
		log.Printf("Successfully updated %s", sub.Name)
	}
}
