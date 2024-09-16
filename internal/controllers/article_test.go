package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proxy/internal/serializer"
	httpTest "proxy/testing"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"proxy/api/posts"
	"proxy/internal/models"
)

func TestGetArticlesByAuthor(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	testServer := httpTest.NewServer()
	defer testServer.Teardown()

	authorID := 1
	articles := []models.Article{
		{ID: 1, Title: "Test Article 1", Body: "Content 1", UserID: authorID},
		{ID: 2, Title: "Test Article 2", Body: "Content 2", UserID: authorID},
	}
	serializedArticles, _ := serializer.SerializeToString(articles)
	mock.ExpectGet("articles:1").SetVal(serializedArticles)

	testServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(articles)
	})

	pc := &posts.Client{
		HttpClient: &http.Client{},
		PostsUrl:   testServer.URL,
	}

	controller := NewArticleController(rdb, pc)
	req, err := http.NewRequest("GET", "/articles?author=1", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	controller.GetArticles(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mock.ExpectGet("author:123")
}

func TestGetArticle(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	testServer := httpTest.NewServer()
	defer testServer.Teardown()

	article := models.Article{
		ID: 1, Title: "Test Article 1", Body: "Content 1", UserID: 1,
	}
	serializedArticles, _ := serializer.SerializeToString(article)
	mock.ExpectGet("articles:1").SetVal(serializedArticles)

	testServer.Mux.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(article)
	})

	pc := &posts.Client{
		HttpClient: &http.Client{},
		PostsUrl:   testServer.URL + "/1",
	}

	controller := NewArticleController(rdb, pc)
	req, err := http.NewRequest("GET", "/articles/1", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	controller.GetArticles(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mock.ExpectGet("author:123")
}
