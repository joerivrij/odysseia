package service

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"net/http"
	"net/url"
)

type PtolemaiosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	secretEndpoint string = "/ptolemaios/v1/secret"
)

func NewPtolemaiosConfig(scheme, baseUrl string, client HttpClient) (*PtolemaiosImpl, error) {
	return &PtolemaiosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (p *PtolemaiosImpl) GetSecret() (*models.ElasticConfigVault, error) {
	path := url.URL{
		Scheme: p.Scheme,
		Host:   p.BaseUrl,
		Path:   secretEndpoint,
	}

	response, err := p.Client.Get(&path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var secret models.ElasticConfigVault
	err = json.NewDecoder(response.Body).Decode(&secret)
	if err != nil {
		return nil, err
	}

	glg.Debug(secret)

	return &secret, nil
}
