package client

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"net/http"
	"net/url"
	"path"
)

type AlexandrosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	alexandrosService string = "alexandros"
	queryWord         string = "search"
)

func NewAlexandrosImpl(scheme, baseUrl string, client HttpClient) (*AlexandrosImpl, error) {
	return &AlexandrosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (a *AlexandrosImpl) Health() (*models.Health, error) {
	healthPath := url.URL{
		Scheme: a.Scheme,
		Host:   a.BaseUrl,
		Path:   path.Join(alexandrosService, version, healthEndPoint),
	}

	return Health(healthPath, a.Client)
}

func (a *AlexandrosImpl) QueryWord(word string) ([]models.Meros, error) {
	urlPath := url.URL{
		Scheme: a.Scheme,
		Host:   a.BaseUrl,
		Path:   path.Join(alexandrosService, version, queryWord),
	}

	q := urlPath.Query()
	q.Set(queryTermWord, word)
	urlPath.RawQuery = q.Encode()

	response, err := a.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var searchModel []models.Meros
	err = json.NewDecoder(response.Body).Decode(&searchModel)
	if err != nil {
		return nil, err
	}

	return searchModel, nil
}
