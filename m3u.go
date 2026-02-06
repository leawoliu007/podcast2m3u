package main

import (
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
)

func M3u(feed *gofeed.Feed, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "#EXTM3U")
	for _, item := range feed.Items {
		url := ""
		if len(item.Enclosures) > 0 {
			url = item.Enclosures[0].URL
		}
		if url == "" {
			continue // Skip items without audio
		}
		
		fmt.Fprintf(file, "#EXTINF:%d,%s - %s\n%s\n",
			-1,
			feed.Title,
			item.Title,
			url)
	}
	return nil
}
