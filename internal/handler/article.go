package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"proxy/api"
	"proxy/internal/config"
	"proxy/internal/serializer"
	"proxy/internal/util"
	"proxy/models"
	"strconv"
)

func ArticlesHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("author") {
		authorIdStr := r.URL.Query().Get("author")
		if authorIdStr == "" {
			util.Error(w, http.StatusBadRequest, "Missing author ID")
			return
		}
		authorId, err := strconv.Atoi(authorIdStr)
		if err != nil {
			util.Error(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		articles, err := api.FetchArticles(&authorId)
		if err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.JSON(w, http.StatusOK, articles)
	}

	hasUuid := r.URL.Query().Has("uuid")
	if hasUuid {
		uuid := r.URL.Query().Get("uuid")
		if uuid == "" {
			util.Error(w, http.StatusNotFound, "Missing uuid")
		}
		viewedPagesIds, _ := rdb.SMembers(context.TODO(), fmt.Sprintf("user:%s:articles", uuid)).Result()
		articles := make([]models.Article, len(viewedPagesIds))
		for i, viewedPageId := range viewedPagesIds {
			id, _ := strconv.Atoi(viewedPageId)
			article, _ := api.FetchArticle(id)
			articles[i] = *article
		}
		util.JSON(w, http.StatusOK, articles)
	}

}

func ArticleHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	article, found := getCachedArticle(rdb, id)
	if found {
		util.JSON(w, http.StatusOK, article)
		return
	}

	article, err := api.FetchArticle(id)
	if err == nil {
		util.JSON(w, http.StatusOK, article)
	}

	util.Error(w, http.StatusInternalServerError, "Failed to fetch articles")
}

func getCachedArticle(rdb *redis.Client, id int) (*models.Article, bool) {
	articleStr, err := rdb.Get(context.TODO(), strconv.Itoa(id)).Result()
	if errors.Is(err, redis.Nil) {
		fetchedArticle, err := api.FetchArticle(id)
		if err != nil {
			log.Printf("Failed to fetch article: %v", err)
			return nil, false
		}
		setCachedArticle(rdb, id, *fetchedArticle)
		return fetchedArticle, true
	} else if err != nil {
		log.Println("Error fetching from Redis:", err)
		return nil, false
	}
	article := &models.Article{}
	err = serializer.DeserializeFromString(articleStr, article)
	return article, true
}

func setCachedArticle(rdb *redis.Client, id int, article models.Article) {
	articleString, err := serializer.SerializeToString(article)
	if err != nil {
		log.Printf("Failed to serialize article: %v", err)
	}
	err = rdb.Set(context.TODO(), strconv.Itoa(id), articleString, config.DetailedArticledCache).Err()
	if err != nil {
		log.Println("Error setting Redis cache:", err)
	}
}
