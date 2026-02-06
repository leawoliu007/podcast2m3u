package main

import (
	"fmt"
	"os"

	"github.com/nsoufr/podfeed"
)

func M3u(podcast podfeed.Podcast, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "#EXTM3U")
	for _, episode := range podcast.Items {
		fmt.Fprintf(file, "#EXTINF:%d,%s - %s\n%s\n",
			-1,
			podcast.Title,
			episode.Title,
			episode.Enclosure.Url)
	}
	return nil
}
