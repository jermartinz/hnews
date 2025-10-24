// Package api manage HN API
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *APIClient {
	return &APIClient{
		baseURL: "https://hacker-news.firebaseio.com/v0",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *APIClient) GetTopStories() ([]int, error) {
	url := fmt.Sprintf("%s/topstories.json", a.baseURL)

	resp, err := a.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching top stories: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing body: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return ids, nil
}
