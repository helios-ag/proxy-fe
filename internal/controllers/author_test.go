package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proxy/api/authors"
	"proxy/internal/serializer"
	httpTest "proxy/testing"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"proxy/internal/models"
)

func TestGetAuthors(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	testServer := httpTest.NewServer()
	defer testServer.Teardown()

	authorsArr := []models.Author{{
		ID:       1,
		Name:     "Wes Borland",
		Username: "wes",
		Email:    "wes@lz.com",
		Phone:    "+123888",
		Website:  "https://google.com",
	}}
	serializedArticles, _ := serializer.SerializeToString(authorsArr)
	mock.ExpectGet("authors").SetVal(serializedArticles)

	testServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authorsArr)
	})

	ac := &authors.Client{
		HttpClient: &http.Client{},
		AuthorsUrl: testServer.URL,
	}

	controller := NewAuthorArticlesController(rdb, ac)
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	controller.GetAuthorArticles(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mock.ExpectGet("authors")
}
