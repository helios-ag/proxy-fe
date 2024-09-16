package controllers

import (
	"github.com/go-redis/redismock/v8"
	"net/http"
	"proxy/api/authors"
	"proxy/api/posts"
	"proxy/internal/config"
	"testing"
)

func TestNewControllerFactory(t *testing.T) {
	db, _ := redismock.NewClientMock()
	pc := &posts.Client{
		HttpClient: &http.Client{},
		PostsUrl:   config.PostsURL,
	}
	ac := &authors.Client{
		HttpClient: &http.Client{},
		AuthorsUrl: config.UsersURL,
	}
	controllerFactory := NewControllerFactory(db, pc, ac)

	if controllerFactory.rdb != db || controllerFactory.postsClient != pc || controllerFactory.authorsClient != ac {
		t.Errorf("NewControllerFactory did not initialize with the correct values")
	}
}

// Test for GetArticleController
func TestGetArticleController(t *testing.T) {
	db, _ := redismock.NewClientMock()
	pc := &posts.Client{
		HttpClient: &http.Client{},
		PostsUrl:   config.PostsURL,
	}
	ac := &authors.Client{
		HttpClient: &http.Client{},
		AuthorsUrl: config.UsersURL,
	}

	controllerFactory := &ControllerFactory{
		rdb:           db,
		postsClient:   pc,
		authorsClient: ac,
	}

	articleController := controllerFactory.GetArticleController()
	if articleController == nil {
		t.Errorf("GetArticleController returned nil")
	}
}

func TestGetAuthorArticlesController(t *testing.T) {
	db, _ := redismock.NewClientMock()

	pc := &posts.Client{
		HttpClient: &http.Client{},
		PostsUrl:   config.PostsURL,
	}
	ac := &authors.Client{
		HttpClient: &http.Client{},
		AuthorsUrl: config.UsersURL,
	}

	controllerFactory := &ControllerFactory{
		rdb:           db,
		authorsClient: ac,
		postsClient:   pc,
	}

	authorArticlesController := controllerFactory.GetAuthorArticlesController()
	if authorArticlesController == nil {
		t.Errorf("GetAuthorArticlesController returned nil")
	}
}

func TestGetTrackController(t *testing.T) {
	db, _ := redismock.NewClientMock()
	controllerFactory := &ControllerFactory{
		rdb: db,
	}

	trackController := controllerFactory.GetTrackController()
	if trackController == nil {
		t.Errorf("GetTrackController returned nil")
	}
}
