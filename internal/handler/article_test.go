package handler_test

import (
	"net/http"
	"net/http/httptest"
	"proxy/internal/handler"
	"proxy/internal/models"
	"proxy/internal/serializer"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestArticlesHandler(t *testing.T) {
	var (
		mockRedisClient *redis.Client
		mock            redismock.ClientMock
		handlerFunc     http.HandlerFunc
	)

	mockRedisClient, mock = redismock.NewClientMock()
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		handler.ArticlesHandler(mockRedisClient, w, r)
	}

	t.Run("should return bad request for missing author ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/articles?author=", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handlerFunc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Missing author ID")
	})

	t.Run("should return bad request for invalid author ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/articles?author=abc", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handlerFunc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid author ID")
	})

	t.Run("should return articles for valid author ID", func(t *testing.T) {
		authorID := 1
		articles := []models.Article{{ID: 1, Title: "Test Article"}}
		serializedArticles, _ := serializer.SerializeToString(articles)

		mock.ExpectGet("articles:" + strconv.Itoa(authorID)).SetVal(serializedArticles)

		req, err := http.NewRequest("GET", "/articles?author=1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handlerFunc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Test Article")
	})
}

func TestArticleHandler(t *testing.T) {
	var (
		mockRedisClient *redis.Client
		mock            redismock.ClientMock
		handlerFunc     http.HandlerFunc
	)

	mockRedisClient, mock = redismock.NewClientMock()
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		handler.ArticleHandler(mockRedisClient, w, r)
	}

	t.Run("should return cached article if found", func(t *testing.T) {
		articleID := 1
		article := models.Article{ID: articleID, Title: "Test Article"}
		serializedArticle, _ := serializer.SerializeToString(article)

		mock.ExpectGet(strconv.Itoa(articleID)).SetVal(serializedArticle)

		req, err := http.NewRequest("GET", "/articles/1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/articles/{id}", handlerFunc)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Test Article")
	})

	t.Run("should return internal server error if cached article not found", func(t *testing.T) {
		articleID := 1

		mock.ExpectGet(strconv.Itoa(articleID)).RedisNil()

		req, err := http.NewRequest("GET", "/articles/1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/articles/{id}", handlerFunc)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Failed to fetch articles")
	})
}
