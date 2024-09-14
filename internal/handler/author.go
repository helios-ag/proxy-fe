package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"proxy/api"
	"proxy/internal/config"
	"proxy/internal/models"
	"proxy/internal/serializer"
	"proxy/internal/util"
)

func AuthorArticlesHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	authors := getCachedAuthorsList(rdb)
	if authors != nil {
		util.JSON(w, http.StatusOK, authors)
		return
	}
	util.Error(w, http.StatusInternalServerError, "Error fetching authors")
}

func setCachedAuthors(rdb *redis.Client, authors []models.Author) {
	serializedAuthors, err := serializer.SerializeToString(authors)
	if err != nil {
		log.Printf("Failed to serialize authors: %v", err)
	}
	err = rdb.Set(context.TODO(), "authors", serializedAuthors, config.CacheAuthorsList).Err()
	if err != nil {
		log.Println("Error setting Redis cache:", err)
	}
}

func getCachedAuthorsList(rdb *redis.Client) []models.Author {
	authorsString, err := rdb.Get(context.TODO(), "authors").Result()
	if errors.Is(err, redis.Nil) {
		fetchedAuthors, err := api.FetchAuthors()
		if err != nil {
			log.Printf("Failed to fetch article: %v", err)
			return nil
		}
		setCachedAuthors(rdb, fetchedAuthors)
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
