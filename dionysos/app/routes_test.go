//go:build !integration
// +build !integration

package app

import (
	"encoding/json"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	dictionaryIndexDefault = "dictionary"
	elasticIndexDefault    = "grammar"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := configs.DionysosConfig{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/dionysos/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealthEndpointHealthy(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.DionysosConfig{
		ElasticClient: *mockElasticClient,
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/health")

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
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := configs.DionysosConfig{
		ElasticClient: *mockElasticClient,
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.False(t, healthModel.Healthy)
}

func TestQueryWordEndpointHappyPathFemFirst(t *testing.T) {
	fixtureFile := "dionysosFemaleNoun"
	mockCode := 200
	expected := "noun - sing - fem - nom"
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	testConfig := configs.DionysosConfig{
		ElasticClient:    *mockElasticClient,
		DictionaryIndex:  dictionaryIndexDefault,
		Index:            elasticIndexDefault,
		DeclensionConfig: *declensionConfig,
	}
	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=μάχη")

	var declensions models.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&declensions)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, len(declensions.Results) == 1)
	assert.Equal(t, expected, declensions.Results[0].Rule)
}

func TestQueryWordEndpointHappyPathMascSecond(t *testing.T) {
	fixtureFile := "dionysosMascNoun"
	mockCode := 200
	expected := "noun - plural - masc - nom"
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	testConfig := configs.DionysosConfig{
		ElasticClient:    *mockElasticClient,
		DictionaryIndex:  dictionaryIndexDefault,
		Index:            elasticIndexDefault,
		DeclensionConfig: *declensionConfig,
	}
	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=πόλεμοι")

	var declensions models.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&declensions)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, len(declensions.Results) == 1)
	assert.Equal(t, expected, declensions.Results[0].Rule)
}

func TestSearchEndPointElasticNoResults(t *testing.T) {
	expected := "no options found"

	elasticClient, err := elastic.CreateElasticClient("test", "test", []string{"http://localhost:9200"})
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	testConfig := configs.DionysosConfig{
		ElasticClient:    *elasticClient,
		DictionaryIndex:  dictionaryIndexDefault,
		Index:            elasticIndexDefault,
		DeclensionConfig: *declensionConfig,
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=πόλεμοι")

	var notFoundError models.NotFoundError
	err = json.NewDecoder(response.Body).Decode(&notFoundError)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, expected, notFoundError.Message.Reason)
}

func TestSearchEndPointPrespositionFound(t *testing.T) {
	fixtureFile := "dionysosPreposition"
	mockCode := 200
	expected := "preposition"
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	testConfig := configs.DionysosConfig{
		ElasticClient:    *mockElasticClient,
		DictionaryIndex:  dictionaryIndexDefault,
		Index:            elasticIndexDefault,
		DeclensionConfig: *declensionConfig,
	}
	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=εἰς")

	var declensions models.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&declensions)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	for _, decl := range declensions.Results {
		assert.Equal(t, expected, decl.Rule)
	}
}

func TestSearchEndPointWithoutQueryParam(t *testing.T) {
	expected := "cannot be empty"

	elasticClient, err := elastic.CreateElasticClient("test", "test", []string{"http://localhost:9200"})
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	testConfig := configs.DionysosConfig{
		ElasticClient:    *elasticClient,
		DictionaryIndex:  dictionaryIndexDefault,
		Index:            elasticIndexDefault,
		DeclensionConfig: *declensionConfig,
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=")

	var validation models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&validation)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, expected, validation.Messages[0].Message)
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
