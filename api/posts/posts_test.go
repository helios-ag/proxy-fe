package posts_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proxy/api/posts"
	"proxy/internal/models"
	"strconv"
	"testing"
)

func TestFetchArticles(t *testing.T) {
	mockArticles := []models.Article{
		{
			UserID: 0,
			ID:     1,
			Title:  "Article 1",
			Body:   "",
			Viewed: false,
		},
		{
			UserID: 0,
			ID:     2,
			Title:  "Article 2",
			Body:   "",
			Viewed: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockArticles)
	}))
	defer server.Close()

	client := posts.Client{
		HttpClient: server.Client(),
		PostsUrl:   server.URL,
	}

	articles, err := client.FetchArticles()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(articles) != len(mockArticles) {
		t.Fatalf("expected %d articles, got %d", len(mockArticles), len(articles))
	}
}

func TestFetchArticlesByAuthor(t *testing.T) {
	authorId := 1

	mockArticles := []models.Article{
		{
			UserID: authorId,
			ID:     1,
			Title:  "Article 1",
			Body:   "",
			Viewed: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("userId") != strconv.Itoa(authorId) {
			t.Fatalf("expected userId %d, got %s", authorId, r.URL.Query().Get("userId"))
		}
		json.NewEncoder(w).Encode(mockArticles)
	}))
	defer server.Close()

	client := posts.Client{
		HttpClient: server.Client(),
		PostsUrl:   server.URL,
	}

	articles, err := client.FetchArticlesByAuthor(authorId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(articles) != len(mockArticles) {
		t.Fatalf("expected %d articles, got %d", len(mockArticles), len(articles))
	}
}

func TestFetchArticle(t *testing.T) {
	articleId := 1
	mockArticle := models.Article{ID: articleId, Title: "Article 1"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+strconv.Itoa(articleId) {
			t.Fatalf("expected path /%d, got %s", articleId, r.URL.Path)
		}
		json.NewEncoder(w).Encode(mockArticle)
	}))
	defer server.Close()

	client := posts.Client{
		HttpClient: server.Client(),
		PostsUrl:   server.URL,
	}

	article, err := client.FetchArticle(articleId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if article.ID != mockArticle.ID {
		t.Fatalf("expected article ID %d, got %d", mockArticle.ID, article.ID)
	}
}
