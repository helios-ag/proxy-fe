package controllers

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"proxy/api/authors"
	"proxy/internal/config"
	"proxy/internal/models"
	"proxy/internal/serializer"
	"proxy/internal/util"
)

func NewAuthorArticlesController(rdb *redis.Client, client *authors.Client) *AuthorArticlesController {
	return &AuthorArticlesController{rdb, client}
}

type AuthorArticlesController struct {
	rdb *redis.Client
	c   *authors.Client
}

func (controller AuthorArticlesController) GetAuthorArticles(w http.ResponseWriter, r *http.Request) {
	authorsList := controller.getCachedAuthorsList(controller.rdb)
	if authorsList != nil {
		util.JSON(w, http.StatusOK, authorsList)
		return
	}
	util.Error(w, http.StatusInternalServerError, "Error fetching authorsList")
}

func (controller AuthorArticlesController) setCachedAuthors(rdb *redis.Client, authors []models.Author) {
	serializedAuthors, err := serializer.SerializeToString(authors)
	if err != nil {
		log.Printf("Failed to serialize authors: %v", err)
	}
	err = rdb.Set(context.TODO(), "authors", serializedAuthors, config.CacheAuthorsList).Err()
	if err != nil {
		log.Println("Error setting Redis cache:", err)
	}
}

func (controller AuthorArticlesController) getCachedAuthorsList(rdb *redis.Client) []models.Author {
	authorsString, err := rdb.Get(context.TODO(), "authors").Result()
	if errors.Is(err, redis.Nil) {
		fetchedAuthors, err := controller.c.FetchAuthors()
		if err != nil {
			log.Printf("Failed to fetch article: %v", err)
			return nil
		}
		controller.setCachedAuthors(rdb, fetchedAuthors)
		return fetchedAuthors
	} else if err != nil {
		log.Println("Error fetching from Redis:", err)
		return nil
	}
	var deserializedAuthorList []models.Author
	err = serializer.DeserializeFromString(authorsString, &deserializedAuthorList)
	if err != nil {
		log.Printf("Failed to deserialize authors: %v", err)
		return nil
	}
	return deserializedAuthorList
}
