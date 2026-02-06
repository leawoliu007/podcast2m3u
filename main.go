package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"github.com/nsoufr/podfeed"
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
		podcast, err := podfeed.Fetch(context.Background(), legacyPodcastURI)
		if err != nil {
			log.Fatal(err)
		}
		// In legacy mode, we print to stdout, so we need a temporary way or just print it.
		// Since M3u now takes a path, we can't easily use it for stdout without modification 
		// or passing /dev/stdout (linux) or similar. 
		// For cross-platform, let's just implement a quick stdout writer or modify M3u slightly later.
		// actually, let's keep it simple: if M3u function takes a filename, we can pass "/dev/stdout" on linux, 
		// but on windows that's tricky.
		// Let's just re-implement the simple print for legacy to avoid breaking it too much, 
		// or better: let's treat legacy as a separate quick flow.
		fmt.Println("#EXTM3U")
		for _, episode := range podcast.Items {
			fmt.Printf("#EXTINF:%d,%s - %s\n%s\n", -1, podcast.Title, episode.Title, episode.Enclosure.Url)
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
