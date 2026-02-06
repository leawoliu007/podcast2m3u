package main

import (
	"fmt"

	"github.com/nsoufr/podfeed"
)

func M3u(podcast podfeed.Podcast) {
	fmt.Println("#EXTM3U")
	for _, episode := range podcast.Items {
		fmt.Printf("#EXTINF:%d,%s - %s\n%s\n",
			-1,
			podcast.Title,
			episode.Title,
			episode.Enclosure.Url)
	}

}
