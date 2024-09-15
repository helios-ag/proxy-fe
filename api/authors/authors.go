package authors

import (
	"encoding/json"
	"io"
	"net/http"
	"proxy/internal/models"
)

type Client struct {
	HttpClient *http.Client
	AuthorsUrl string
}

func (c Client) FetchAuthors() ([]models.Author, error) {
	resp, err := c.HttpClient.Get(c.AuthorsUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authors []models.Author
	err = json.Unmarshal(body, &authors)
	if err != nil {
		return nil, err
	}

	return authors, nil
}
