package client

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"net/http"
	"net/url"
	"path"
)

type HerodotosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	herodotosService string = "herodotos"
	checkSentence    string = "checkSentence"
	queryTermAuthor  string = "author"
	queryTermBook    string = "book"
	authors          string = "authors"
	books            string = "books"
)

func NewHerodotosImpl(scheme, baseUrl string, client HttpClient) (*HerodotosImpl, error) {
	return &HerodotosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (h *HerodotosImpl) Health() (*models.Health, error) {
	healthPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   path.Join(herodotosService, version, healthEndPoint),
	}

	return Health(healthPath, h.Client)
}

func (h *HerodotosImpl) Authors() (*models.Authors, error) {
	urlPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   path.Join(herodotosService, version, authors),
	}

	response, err := h.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	var authorModel models.Authors
	err = json.NewDecoder(response.Body).Decode(&authorModel)
	if err != nil {
		return nil, err
	}

	return &authorModel, nil
}

func (h *HerodotosImpl) Books(author string) (*models.Books, error) {
	urlPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   path.Join(herodotosService, version, authors, author, books),
	}

	response, err := h.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	var booksModel models.Books
	err = json.NewDecoder(response.Body).Decode(&booksModel)
	if err != nil {
		return nil, err
	}

	return &booksModel, nil
}

func (h *HerodotosImpl) CreateQuestion(author, book string) (*models.CreateSentenceResponse, error) {
	urlPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   path.Join(herodotosService, version, createQuestion),
	}

	q := urlPath.Query()
	q.Set(queryTermAuthor, author)
	q.Set(queryTermBook, book)
	urlPath.RawQuery = q.Encode()

	response, err := h.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var sentenceResponse models.CreateSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&sentenceResponse)
	if err != nil {
		return nil, err
	}

	return &sentenceResponse, nil
}

func (h *HerodotosImpl) CheckSentence(requestBody models.CheckSentenceRequest) (*models.CheckSentenceResponse, error) {
	urlPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   path.Join(herodotosService, version, checkSentence),
	}

	body, err := requestBody.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := h.Client.Post(&urlPath, body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var sentenceResponse models.CheckSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&sentenceResponse)
	if err != nil {
		return nil, err
	}

	return &sentenceResponse, nil
}
