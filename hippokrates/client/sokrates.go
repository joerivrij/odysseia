package client

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"net/http"
	"net/url"
	"path"
)

type SokratesImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	sokratesService string = "sokrates"
	methods         string = "methods"
	categories      string = "categories"
	chapters        string = "chapters"
	queryMethod     string = "method"
	queryCategory   string = "category"
	queryChapter    string = "chapter"
	answer          string = "answer"
)

func NewSokratesImpl(scheme, baseUrl string, client HttpClient) (*SokratesImpl, error) {
	return &SokratesImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (s *SokratesImpl) Health() (*models.Health, error) {
	healthPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, healthEndPoint),
	}

	return Health(healthPath, s.Client)
}

func (s *SokratesImpl) Methods() (*models.Methods, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, methods),
	}

	response, err := s.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	var methodModel models.Methods
	err = json.NewDecoder(response.Body).Decode(&methodModel)
	if err != nil {
		return nil, err
	}

	return &methodModel, nil
}

func (s *SokratesImpl) Categories(method string) (*models.Categories, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, methods, method, categories),
	}

	response, err := s.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	var categoryModel models.Categories
	err = json.NewDecoder(response.Body).Decode(&categoryModel)
	if err != nil {
		return nil, err
	}

	return &categoryModel, nil
}

func (s *SokratesImpl) LastChapter(method, category string) (*models.LastChapterResponse, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, methods, method, categories, category, chapters),
	}

	response, err := s.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	var lastChapter models.LastChapterResponse
	err = json.NewDecoder(response.Body).Decode(&lastChapter)
	if err != nil {
		return nil, err
	}

	return &lastChapter, nil
}

func (s *SokratesImpl) CreateQuestion(method, category, chapter string) (models.QuizResponse, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, createQuestion),
	}

	q := urlPath.Query()
	q.Set(queryMethod, method)
	q.Set(queryCategory, category)
	q.Set(queryChapter, chapter)
	urlPath.RawQuery = q.Encode()

	response, err := s.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var quizResponse models.QuizResponse
	err = json.NewDecoder(response.Body).Decode(&quizResponse)
	if err != nil {
		return nil, err
	}

	return quizResponse, nil
}

func (s *SokratesImpl) Answer(request models.CheckAnswerRequest) (*models.CheckAnswerResponse, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, answer),
	}

	body, err := request.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := s.Client.Post(&urlPath, body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	var checkAnswerModel models.CheckAnswerResponse
	err = json.NewDecoder(response.Body).Decode(&checkAnswerModel)
	if err != nil {
		return nil, err
	}

	return &checkAnswerModel, nil
}
