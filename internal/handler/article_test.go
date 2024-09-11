package handler

import (
	"github.com/go-redis/redismock/v8"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticlesHandler(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Mock the SMembers call
	mock.ExpectSMembers("articles").SetVal([]string{"1"})

	req, err := http.NewRequest("GET", "/articles/1", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ArticlesHandler(db, w, r)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestArticleHandler(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Mock the Get call
	mock.ExpectGet("article:1").RedisNil()

	req, err := http.NewRequest("GET", "/article?id=1", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ArticleHandler(db, w, r)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
