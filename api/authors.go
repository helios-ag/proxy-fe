package api

import (
	"encoding/json"
	"io"
	"net/http"
	"proxy/internal/config"
	"proxy/internal/models"
)

func FetchAuthors() ([]models.Author, error) {
	resp, err := http.Get(config.UsersURL)
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
