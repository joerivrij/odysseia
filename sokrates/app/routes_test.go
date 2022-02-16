package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := configs.SokratesConfig{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/sokrates/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealthEndpointHealthy(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/sokrates/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	//models.Health
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, healthModel.Healthy)
}

func TestHealthEndpointElasticDown(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/sokrates/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.False(t, healthModel.Healthy)
}

func TestLastChapterEndPointHappyPath(t *testing.T) {
	fixtureFile := "lastChapterSokrates"
	mockCode := 200
	expectedChapter := int64(15)
	category := "nomina"
	method := "mousieon"
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/methods/%s/categories/%s/chapters", method, category))

	var searchResults models.LastChapterResponse
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, expectedChapter, searchResults.LastChapter)
}

func TestLastChapterShardFailure(t *testing.T) {
	fixtureFile := "infoServiceDown"
	mockCode := 502
	category := "nomina"
	method := "method"
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/methods/%s/categories/%s/chapters", method, category))

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "elasticSearch returned an error"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestLastChapterEmptyQuery(t *testing.T) {
	fixtureFile := "lastChapterSokrates"
	mockCode := 200
	category := "f"
	method := "s"
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/methods/%s/categories/%s/chapters", method, category))

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "must be longer than 1", searchResults.Messages[0].Message)
	assert.Equal(t, "category", searchResults.Messages[0].Field)
}

func TestCheckAnswerEndPointHappyPath(t *testing.T) {
	fixtureFile := "checkQuestionSokrates"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	body := map[string]string{"answerProvided": "godin", "quizWord": "θεός", "category": "nomina"}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/sokrates/v1/answer", bodyInBytes)

	var searchResults models.CheckAnswerResponse
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, searchResults.Correct)
}

func TestCheckAnswerElasticDown(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	body := map[string]string{"answerProvided": "godin", "quizWord": "θεός", "category": "nomina"}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/sokrates/v1/answer", bodyInBytes)

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "elasticSearch returned an error"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestCheckQuestionBadJson(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	body := map[int]string{1: "34"}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/sokrates/v1/answer", bodyInBytes)

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "error parsing", searchResults.Messages[0].Message)
	assert.Contains(t, searchResults.Messages[0].Field, "body")
}

func TestCreateQuestionEndPointHappyPath(t *testing.T) {
	fixtureFile := "createQuestionSokrates"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	category := "verba"
	chapter := "1"
	method := "logos"

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?method=%s&category=%s&chapter=%s", method, category, chapter))

	var searchResults models.QuizResponse
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)

	expectedGreek := [5]string{"ἀναβαίνω", "λέγω", "προβαίνω", "πονέω", "φέπω"}

	assert.Contains(t, expectedGreek, searchResults[0])
}

func TestCreateQuestionEndPointCanCreateShorterQuiz(t *testing.T) {
	fixtureFile := "createQuestionSokratesShort"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	category := "verba"
	method := "mousieon"
	chapter := "1"

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?method=%s&category=%s&chapter=%s", method, category, chapter))

	var searchResults models.QuizResponse
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)

	expectedGreek := [5]string{"ἀναβαίνω", "λέγω", "προβαίνω", "πονέω", "φέπω"}

	assert.Contains(t, expectedGreek, searchResults[0])
}

func TestCreateQuestionEmptyQuery(t *testing.T) {
	fixtureFile := "createQuestionSokratesShort"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "test",
	}

	category := "verba"
	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?category=%s", category))

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "cannot be empty", searchResults.Messages[0].Message)
	assert.Contains(t, searchResults.Messages[0].Field, "chapter")
}

func TestCreateQuestionElasticDown(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	category := "verba"
	chapter := "1"

	testConfig := configs.SokratesConfig{
		ElasticClient: *mockElasticClient,
		SearchWord:    "greek",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?category=%s&chapter=%s", category, chapter))

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	expectedText := "elasticSearch returned an error"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
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
