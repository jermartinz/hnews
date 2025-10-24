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

	for _, item := range items {
		fmt.Printf("ID: %d\nTitle: %s\nURL: %s\n", item.ID, item.Title, item.URL)
	}
}
