package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"net/http"
	"crypto/tls"
	"github.com/robfig/cron/v3"
	"path/filepath"
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
	var podcast *podfeed.Podcast
	var err error

	if globalConfig.SkipCertVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		resp, err := client.Get(sub.URL)
		if err != nil {
			log.Printf("Failed to fetch podcast (insecure) %s: %v", sub.Name, err)
			return
		}
		defer resp.Body.Close()

		podcast, err = podfeed.Parse(resp.Body)
	} else {
		podcast, err = podfeed.Fetch(context.Background(), sub.URL)
	}

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
	err = M3u(podcast, outputPath)
	if err != nil {
		log.Printf("Failed to write M3U for %s: %v", sub.Name, err)
	} else {
		log.Printf("Successfully updated %s", sub.Name)
	}
}
