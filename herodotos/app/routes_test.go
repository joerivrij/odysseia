package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := configs.HerodotosConfig{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/herodotos/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealthEndpointHealthy(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	//models.Health
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, healthModel.Healthy)
}

func TestHealthEndpointElasticDown(t *testing.T) {
	fixtureFile := "infoServiceDown"
	mockCode := 502
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.False(t, healthModel.Healthy)
}

func TestAuthorsEndPointHappyPath(t *testing.T) {
	fixtureFile := "authors"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/authors")

	var searchResults models.Authors
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 3, len(searchResults.Authors))

	expectedAuthors := [3]string{"herodotos", "ploutarchos", "thucydides"}

	for _, author := range searchResults.Authors {
		assert.Contains(t, expectedAuthors, author.Author)
	}
}

func TestAuthorsEndPointShardFailure(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/authors")

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "informations from Elasticsearch"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestBooksEndPointHappyPath(t *testing.T) {
	fixtureFile := "herodotosBooks"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/authors/testAuthor/books")

	var searchResults models.Books
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 1, len(searchResults.Books))

	for _, expectedKey := range searchResults.Books {
		assert.Equal(t, 1, int(expectedKey.Book))
	}
}

func TestBookEndPointShardFailure(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/authors/testAuthor/books")

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "informations from Elasticsearch"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestCreateQuestionHappyPath(t *testing.T) {
	fixtureFile := "createQuestionHerodotos"
	mockCode := 200
	author := "thucydides"
	book := 1
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/herodotos/v1/createQuestion?author=%s&book=%v", author, book))

	var searchResults models.CreateSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "Θουκυδίδης"

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, searchResults.Sentence, expectedText)
}

func TestCreateQuestionMissingAuthor(t *testing.T) {
	fixtureFile := "createQuestionHerodotos"
	mockCode := 200
	author := ""
	book := 1
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/herodotos/v1/createQuestion?author=%s&book=%v", author, book))

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "cannot be empty"

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, searchResults.Messages[0].Message, expectedText)
}

func TestCreateQuestionMissingBook(t *testing.T) {
	fixtureFile := "createQuestionHerodotos"
	mockCode := 200
	author := "someauthor"
	book := ""
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/herodotos/v1/createQuestion?author=%s&book=%s", author, book))

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "cannot be empty"

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, searchResults.Messages[0].Message, expectedText)
}

func TestCreateNoResults(t *testing.T) {
	fixtureFile := "searchWordNoResults"
	mockCode := 200
	author := "someauthor"
	book := "1"
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/herodotos/v1/createQuestion?author=%s&book=%s", author, book))

	var searchResults models.NotFoundError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "no hits for combination"

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, searchResults.Message.Reason, expectedText)
}

func TestCreateQuestionMissingAuthorInElastic(t *testing.T) {
	fixtureFile := "createQuestionHerodotos"
	mockCode := 404
	book := 1
	author := "notanauthor"
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	expectedText := "informations from Elasticsearch"

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/herodotos/v1/createQuestion?author=%s&book=%v", author, book))

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestCreateQuestionShardFailure(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	book := 1
	author := "someAuthor"
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/herodotos/v1/createQuestion?author=%s&book=%v", author, book))

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "informations from Elasticsearch"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestCreateQuestionUnParseableJson(t *testing.T) {
	fixtureFile := "withAll"
	mockCode := 200
	author := "thucydides"
	book := "1"
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/herodotos/v1/createQuestion?author=%s&book=%v", author, book))

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := [2]string{"createQuestion", "translation"}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	for _, message := range searchResults.Messages {
		assert.Contains(t, expectedText, message.Field)
	}
}

func TestCheckAnswerEndPointHappyPath(t *testing.T) {
	fixtureFile := "checkSentenceHerodotos"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	body := map[string]string{"answerSentence": "The Foenicians. ;came to Argos,,.;:'' afd set out some cargo",
		"sentenceId": "GmBFYHkBkbwXxxT5S6F_",
		"author":     "herodotos"}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/herodotos/v1/checkSentence", bodyInBytes)

	var searchResults models.CheckSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	levenshteinAsFloat, _ := strconv.ParseFloat(searchResults.LevenshteinPercentage, 32)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, levenshteinAsFloat > 50)
}

func TestCheckAnswerEndPointBadJsonRequest(t *testing.T) {
	fixtureFile := "checkSentenceHerodotos"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	body := []byte{78}

	bodyInBytes := bytes.NewReader(body)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/herodotos/v1/checkSentence", bodyInBytes)

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	expected := "invalid character"
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, searchResults.Messages[0].Message, expected)
}

func TestCheckSentenceShardFailure(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	body := map[string]string{"answerSentence": "The Foenicians. ;came to Argos,,.;:'' afd set out some cargo",
		"sentenceId": "GmBFYHkBkbwXxxT5S6F_",
		"author":     "herodotos"}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/herodotos/v1/checkSentence", bodyInBytes)

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "informations from Elasticsearch"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestCheckSentenceUnparseableJson(t *testing.T) {
	fixtureFile := "withAll"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	body := map[string]string{"answerSentence": "The Foenicians. ;came to Argos,,.;:'' afd set out some cargo",
		"sentenceId": "GmBFYHkBkbwXxxT5S6F_",
		"author":     "herodotos"}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)

	testConfig := configs.HerodotosConfig{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/herodotos/v1/checkSentence", bodyInBytes)
	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := [2]string{"createQuestion", "translation"}

	assert.Equal(t, http.StatusBadRequest, response.Code)
	for _, message := range searchResults.Messages {
		assert.Contains(t, expectedText, message.Field)
	}
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
