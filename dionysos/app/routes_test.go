package app

import (
	"encoding/json"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := DionysosConfig{}
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

	testConfig := DionysosConfig{
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

	testConfig := DionysosConfig{
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

func TestQueryWordEndpointHappyPath(t *testing.T) {
	fixtureFile := "dionysosFemaleNoun"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	testConfig := DionysosConfig{
		ElasticClient:      *mockElasticClient,
		DictionaryIndex: dictionaryIndexDefault,
		Index:             elasticIndexDefault,
		DeclensionConfig:   *declensionConfig,
	}
	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=μάχη")

	var declensions models.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&declensions)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
