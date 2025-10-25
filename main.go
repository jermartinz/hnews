package main

import (
	"fmt"
	"log"

	"github.com/jermartinz/hn/api"
)

func main() {
	client := api.NewClient()

	items, err := client.GetItemStories()
	if err != nil {
		log.Fatal(err)
	}

	for i, item := range items {
		if i >= 10 {
			break
		}
		fmt.Printf("Title: %s\nURL: %s\n", item.Title, item.URL)
	}
}
