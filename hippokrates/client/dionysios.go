package client

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"net/http"
	"net/url"
	"path"
)

type DionysiosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	dionysiosService string = "dionysios"
	checkGrammar     string = "checkGrammar"
)

func NewDionysiosImpl(scheme, baseUrl string, client HttpClient) (*DionysiosImpl, error) {
	return &DionysiosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (d *DionysiosImpl) Health() (*models.Health, error) {
	healthPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(dionysiosService, version, healthEndPoint),
	}

	return Health(healthPath, d.Client)
}

func (d *DionysiosImpl) CheckGrammar(word string) (*models.DeclensionTranslationResults, error) {
	urlPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(dionysiosService, version, checkGrammar),
	}

	q := urlPath.Query()
	q.Set(queryTermWord, word)
	urlPath.RawQuery = q.Encode()

	response, err := d.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var declensionModels models.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&declensionModels)
	if err != nil {
		return nil, err
	}

	return &declensionModels, nil
}
