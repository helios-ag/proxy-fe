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
	"proxy/internal/cookies"
	"proxy/internal/models"
	"proxy/internal/serializer"
	"proxy/internal/util"
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
		articles, err := getAuthorCachedArticles(rdb, authorId)
		if err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.JSON(w, http.StatusOK, articles)
		return
	}

	uuid, err := cookies.Read(r, "userId")
	if uuid == "" {
		util.Error(w, http.StatusNotFound, "Missing uuid")
		return
	}

	if err == nil {
		log.Printf("got uuid: %s", uuid)

		viewedPagesIds, _ := rdb.SMembers(context.TODO(), fmt.Sprintf("user:%s:articles", uuid)).Result()
		if len(viewedPagesIds) == 0 {
			util.Error(w, http.StatusNotFound, "Articles not found")
			return
		}
		articles := make([]models.Article, len(viewedPagesIds))
		for i, viewedPageId := range viewedPagesIds {
			id, _ := strconv.Atoi(viewedPageId)
			article, _ := api.FetchArticle(id)
			articles[i] = *article
		}
		util.JSON(w, http.StatusOK, articles)
		return
	}

	articles, err := api.FetchArticles()
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(w, http.StatusOK, articles)
}

func ArticleHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	article, err := getCachedArticle(rdb, id)
	if err != nil {
		log.Println("Found cached article!")
		util.JSON(w, http.StatusOK, article)
		return
	}

	article, err = api.FetchArticle(id)
	if err == nil {
		util.JSON(w, http.StatusOK, article)
		return
	}

	util.Error(w, http.StatusInternalServerError, "Failed to fetch articles")
}

func getCachedArticle(rdb *redis.Client, id int) (*models.Article, error) {
	articleStr, err := rdb.Get(context.TODO(), strconv.Itoa(id)).Result()
	if errors.Is(err, redis.Nil) {
		fetchedArticle, err := api.FetchArticle(id)
		if err != nil {
			log.Printf("Failed to fetch article: %v", err)
			return nil, err
		}
		setCachedArticle(rdb, id, *fetchedArticle)
		return fetchedArticle, nil
	} else if err != nil {
		log.Println("Error fetching from Redis:", err)
		return nil, err
	}
	article := &models.Article{}
	err = serializer.DeserializeFromString(articleStr, article)
	return article, nil
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

func getAuthorCachedArticles(rdb *redis.Client, authorId int) ([]models.Article, error) {
	articleStrings, err := rdb.Get(context.TODO(), "articles:"+strconv.Itoa(authorId)).Result()
	if errors.Is(err, redis.Nil) {
		fetchedArticles, err := api.FetchArticlesByAuthor(authorId)
		if err != nil {
			log.Printf("Failed to fetch article: %v", err)
			return nil, err
		}
		setAuthorCachedArticle(rdb, authorId, fetchedArticles)
		return fetchedArticles, nil
	} else if err != nil {
		log.Println("Error fetching from Redis:", err)
		return nil, err
	}
	var articles []models.Article
	err = serializer.DeserializeFromString(articleStrings, articles)
	return articles, nil
}

func setAuthorCachedArticle(rdb *redis.Client, id int, articles []models.Article) {
	articlesStrings, err := serializer.SerializeToString(articles)
	if err != nil {
		log.Printf("Failed to serialize articles: %v", err)
	}
	err = rdb.Set(context.TODO(), "articles:"+strconv.Itoa(id), articlesStrings, config.CacheAuthorArticles).Err()
	if err != nil {
		log.Println("Error setting Redis cache:", err)
	}
}
