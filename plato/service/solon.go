package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/odysseia/plato/models"
	"net/http"
	"net/url"
	"path"
)

type SolonImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	solonService     string = "solon"
	tokenEndpoint    string = "token"
	registerEndpoint string = "register"
)

func NewSolonImpl(scheme, baseUrl string, ca []byte, certs []tls.Certificate) (*SolonImpl, error) {
	client := NewHttpClient(ca, certs)
	return &SolonImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func NewFakeSolonImpl(scheme, baseUrl string, client HttpClient) (*SolonImpl, error) {
	return &SolonImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (s *SolonImpl) OneTimeToken() (*models.TokenResponse, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(solonService, version, tokenEndpoint),
	}

	response, err := s.Client.Get(&urlPath)
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
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(solonService, version, registerEndpoint),
	}

	body, err := requestBody.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := s.Client.Post(&urlPath, body)
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
	healthPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(solonService, version, healthEndPoint),
	}

	return Health(healthPath, s.Client)
}
