package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"github.com/mmcdole/gofeed"
)

func main() {
	var configPath string
	var daemonMode bool
	var legacyPodcastURI string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&daemonMode, "daemon", false, "Run in daemon mode with scheduler")
	flag.StringVar(&legacyPodcastURI, "podcast", "", "Legacy: URL of a single podcast to print to stdout")
	flag.Parse()

	// Legacy mode support
	if legacyPodcastURI != "" {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(legacyPodcastURI)
		if err != nil {
			log.Fatal(err)
		}
		// Use WriteM3u with stdout
		err = WriteM3u(feed, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Load Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %v", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Init DB
	InitDB(config.Global.DatabasePath)

	if daemonMode {
		// Start Web Server
		go StartWebServer(config.Global.WebServerPort, configPath, config)

		StartScheduler(config)
	} else {
		// Run once for all subscriptions
		log.Println("Running one-time update for all subscriptions...")
		for _, sub := range config.Subscriptions {
			processSubscription(sub, config.Global)
		}
	}
}
