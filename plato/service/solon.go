package service

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia/plato/models"
	"net/http"
	"net/url"
)

type SolonImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	tokenEndpoint    string = "/solon/v1/token"
	healthEndPoint   string = "/solon/v1/health"
	registerEndpoint string = "/solon/v1/register"
)

func NewSolonImpl(scheme, baseUrl string, client HttpClient) (*SolonImpl, error) {
	return &SolonImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (s *SolonImpl) OneTimeToken() (*models.TokenResponse, error) {
	path := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   tokenEndpoint,
	}

	response, err := s.Client.Get(&path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var tokenModel models.TokenResponse
	err = json.NewDecoder(response.Body).Decode(&tokenModel)
	if err != nil {
		return nil, err
	}

	return &tokenModel, nil
}

func (s *SolonImpl) Register(requestBody models.SolonCreationRequest) (*models.SolonResponse, error) {
	path := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   registerEndpoint,
	}

	body, err := requestBody.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := s.Client.Post(&path, body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var solonResponse models.SolonResponse
	err = json.NewDecoder(response.Body).Decode(&solonResponse)
	if err != nil {
		return nil, err
	}

	return &solonResponse, nil
}

func (s *SolonImpl) Health() (*models.Health, error) {
	path := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   healthEndPoint,
	}

	response, err := s.Client.Get(&path)
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
