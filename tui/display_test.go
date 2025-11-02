package tui

import (
	"errors"
	"strings"
	"testing"

	"github.com/jermartinz/hnews/models"
)

func TestCheckWSl(t *testing.T) {
	content := []byte("Linux version 5.4.72-microsoft-standard")

	if !checkWsl(content) {
		t.Errorf("Expected to detect WSL")
	}
}

func TestModelViewLoading(t *testing.T) {
	test := []struct {
		name         string
		model        Model
		wantContains string
	}{
		{
			name:         "Loading state",
			model:        Model{Loading: true},
			wantContains: "Loading",
		},
		{
			name:         "Error state",
			model:        Model{Err: errors.New("Failed to load")},
			wantContains: "Error: Failed to load",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.View()
			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("View() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestModelViewItems(t *testing.T) {
	m := ModelStyle()
	m.Loading = false
	m.Items = []*models.Item{
		{ItemTitle: "Test Item", URL: "http://example.com/1"},
	}

	got := m.View()

	if strings.Contains(got, "Loading Hacker News...") {
		t.Errorf("View() should not show loading message")
	}

	if strings.Contains(got, "Error:") {
		t.Errorf("View() should not show error message")
	}

	if strings.Contains(got, "No stories available.") {
		t.Errorf("View() should show items")
	}
}
