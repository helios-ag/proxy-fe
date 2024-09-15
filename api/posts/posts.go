package posts

import (
	"encoding/json"
	"io"
	"net/http"
	"proxy/internal/models"
	"strconv"
)

type Client struct {
	HttpClient *http.Client
	PostsUrl   string
}

func (c Client) FetchArticles() ([]models.Article, error) {
	url := c.PostsUrl
	return fetchArticlesFromURL(url, c.HttpClient)
}

func (c Client) FetchArticlesByAuthor(authorId int) ([]models.Article, error) {
	url := c.PostsUrl + "?userId=" + strconv.Itoa(authorId)
	return fetchArticlesFromURL(url, c.HttpClient)
}

func (c Client) FetchArticle(id int) (*models.Article, error) {
	url := c.PostsUrl + "/" + strconv.FormatUint(uint64(id), 10)
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

func fetchArticlesFromURL(url string, http *http.Client) ([]models.Article, error) {
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
