package client

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"net/http"
	"net/url"
	"path"
)

type DionysosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	dionysosService string = "dionysos"
	checkGrammar    string = "checkGrammar"
)

func NewDionysosImpl(scheme, baseUrl string, client HttpClient) (*DionysosImpl, error) {
	return &DionysosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (d *DionysosImpl) Health() (*models.Health, error) {
	healthPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(dionysosService, version, healthEndPoint),
	}

	return Health(healthPath, d.Client)
}

func (d *DionysosImpl) CheckGrammar(word string) (*models.DeclensionTranslationResults, error) {
	urlPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(dionysosService, version, checkGrammar),
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
