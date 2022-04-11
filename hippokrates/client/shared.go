package client

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"net/http"
	"net/url"
)

const (
	version        string = "v1"
	healthEndPoint string = "health"
	queryTermWord  string = "word"
	createQuestion string = "createQuestion"
)

func Health(path url.URL, client HttpClient) (*models.Health, error) {
	response, err := client.Get(&path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var healthResponse models.Health
	err = json.NewDecoder(response.Body).Decode(&healthResponse)
	if err != nil {
		return nil, err
	}

	return &healthResponse, nil
}
