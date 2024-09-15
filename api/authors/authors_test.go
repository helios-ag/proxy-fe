package authors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proxy/internal/models"
	"testing"
)

func TestFetchAuthors(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authors := []models.Author{
			{ID: 1, Name: "Author One"},
			{ID: 2, Name: "Author Two"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authors)
	}))
	defer mockServer.Close()

	client := Client{
		HttpClient: http.DefaultClient,
		AuthorsUrl: mockServer.URL,
	}

	authors, err := client.FetchAuthors()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(authors) != 2 {
		t.Fatalf("Expected 2 authors, got %d", len(authors))
	}

	if authors[0].Name != "Author One" || authors[1].Name != "Author Two" {
		t.Fatalf("Unexpected authors: %v", authors)
	}
}

func TestFetchAuthorsError(t *testing.T) {
	// Create a mock server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	client := Client{
		HttpClient: http.DefaultClient,
		AuthorsUrl: mockServer.URL,
	}

	_, err := client.FetchAuthors()
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
}
