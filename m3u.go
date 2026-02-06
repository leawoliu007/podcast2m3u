package main

import (
	"fmt"
	"io"
	"os"

	"strconv"
	"strings"


	"github.com/mmcdole/gofeed"
)

func M3u(feed *gofeed.Feed, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return WriteM3u(feed, file)
}

func WriteM3u(feed *gofeed.Feed, writer io.Writer) error {
	fmt.Fprintln(writer, "#EXTM3U")
	for _, item := range feed.Items {
		url := ""
		if len(item.Enclosures) > 0 {
			url = item.Enclosures[0].URL
		}
		if url == "" {
			continue // Skip items without audio
		}
		
		// Parse Duration
		duration := -1
		if item.ITunesExt != nil && item.ITunesExt.Duration != "" {
			duration = parseDuration(item.ITunesExt.Duration)
		}

		// Format Title with Date
		title := item.Title
		if item.PublishedParsed != nil {
			title = fmt.Sprintf("[%s] %s", item.PublishedParsed.Format("2006-01-02"), title)
		} else if item.Published != ""{
             // Try to use raw string if parsed is nil, though format might vary
             title = fmt.Sprintf("[%s] %s", item.Published, title)
        }

		fmt.Fprintf(writer, "#EXTINF:%d,%s - %s\n%s\n",
			duration,
			feed.Title,
			title,
			url)
	}
	return nil
}

func parseDuration(d string) int {
	parts := strings.Split(d, ":")
	seconds := 0
	multiplier := 1

	for i := len(parts) - 1; i >= 0; i-- {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			return -1 
		}
		seconds += val * multiplier
		multiplier *= 60
	}
	return seconds
}
