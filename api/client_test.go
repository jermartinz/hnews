package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGetTopStories(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		wantIDs        []int
		wantErr        bool
		errContains    string
	}{
		{
			name:           "succes - returns story IDs",
			serverResponse: `[1, 2, 3, 4, 5]`,
			serverStatus:   http.StatusOK,
			wantIDs:        []int{1, 2, 3, 4, 5},
			wantErr:        false,
		},
		{
			name:           "error - server returns 500",
			serverResponse: ``,
			serverStatus:   http.StatusInternalServerError,
			wantIDs:        nil,
			wantErr:        true,
			errContains:    "unexpected status code",
		},
		{
			name:           "error - invalid JSON",
			serverResponse: `[1, 2, invalid]`,
			serverStatus:   http.StatusOK,
			wantIDs:        nil,
			wantErr:        true,
			errContains:    "error decoding",
		},
		{
			name:           "succes empty array",
			serverResponse: `[]`,
			serverStatus:   http.StatusOK,
			wantIDs:        []int{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer mockServer.Close()

			client := NewClient(mockServer.URL)

			ids, err := client.GetTopStories()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got  %q ", tt.errContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(ids, tt.wantIDs) {
				t.Errorf("got %v, want %v", ids, tt.wantIDs)
			}
		})
	}
}
