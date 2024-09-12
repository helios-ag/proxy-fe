package api

import (
	"encoding/json"
	"io"
	"net/http"
	"proxy/models"
	"strconv"
)

const (
	posts = "https://jsonplaceholder.typicode.com/posts"
)

func FetchArticles() ([]models.Article, error) {
	url := posts

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var articles []models.Article
	err = json.Unmarshal(body, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func FetchArticlesByAuthor(authorId int) (*models.Article, error) {
	url := posts + "?userId=" + strconv.Itoa(authorId)
	return fetchArticleFromURL(url)
}

func FetchArticle(id int) (*models.Article, error) {
	url := posts + "/" + strconv.FormatUint(uint64(id), 10)
	return fetchArticleFromURL(url)
}

func fetchArticleFromURL(url string) (*models.Article, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var article models.Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}
