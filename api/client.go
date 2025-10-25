// Package api manage HN API
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jermartinz/hn/models"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL ...string) *APIClient {
	url := "https://hacker-news.firebaseio.com/v0"
	if len(baseURL) > 0 {
		url = baseURL[0]
	}
	return &APIClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *APIClient) fetch(url string, target any) error {
	resp, err := a.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing body: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

func (a *APIClient) GetTopStories() ([]int, error) {
	url := fmt.Sprintf("%s/topstories.json", a.baseURL)
	var ids []int
	if err := a.fetch(url, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (a *APIClient) GetItemStories() ([]*models.Item, error) {
	storyID, err := a.GetTopStories()
	if err != nil {
		log.Fatalf("error retrieving stories: %v", err)
	}
	var wg sync.WaitGroup
	ch := make(chan *models.Item, len(storyID))
	for _, id := range storyID {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			url := fmt.Sprintf("%s/item/%d.json", a.baseURL, id)
			item := &models.Item{}

			if err := a.fetch(url, item); err != nil {
				return
			}
			ch <- item
		}(id)
	}
	wg.Wait()
	close(ch)

	var items []*models.Item
	for item := range ch {
		items = append(items, item)
	}

	return items, nil
}
