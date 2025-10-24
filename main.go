package main

import (
	"fmt"
	"log"

	"github.com/jermartinz/hn/api"
)

func main() {
	client := api.NewClient()

	ids, err := client.GetTopStories()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ids)
}
