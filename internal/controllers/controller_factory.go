package controllers

import (
	"github.com/go-redis/redis/v8"
	"proxy/api/authors"
	"proxy/api/posts"
)

type ControllerFactory struct {
	rdb           *redis.Client
	postsClient   *posts.Client
	authorsClient *authors.Client
}

func (cf ControllerFactory) GetArticleController() *ArticleController {
	return NewArticleController(cf.rdb, cf.postsClient)
}

func (cf ControllerFactory) GetAuthorArticlesController() *AuthorArticlesController {
	return NewAuthorArticlesController(cf.rdb, cf.authorsClient)
}

func (cf ControllerFactory) GetTrackController() *TrackController {
	return NewTrackController(cf.rdb)
}

func NewControllerFactory(rdb *redis.Client, pc *posts.Client, ac *authors.Client) *ControllerFactory {
	return &ControllerFactory{
		rdb:           rdb,
		postsClient:   pc,
		authorsClient: ac,
	}
}
