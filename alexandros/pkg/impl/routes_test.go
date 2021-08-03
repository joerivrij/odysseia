package impl

import (
	"encoding/json"
	"github.com/odysseia/alexandros/pkg/config"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := config.AlexandrosConfig{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/alexandros/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealthEndpoint(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := config.AlexandrosConfig{
		ElasticClient: *mockElasticClient,
		Index:         "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/alexandros/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	//models.Health
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, healthModel.Healthy)
}

func TestSearchEndPointHappyPath(t *testing.T) {
	fixtureFile := "searchWord"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := config.AlexandrosConfig{
		ElasticClient: *mockElasticClient,
		Index:         "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/alexandros/v1/search?word=αγο")

	var searchResults []models.Meros
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 2, len(searchResults))

	expectedGreek := [2]string{"ἀγορεύω", "ἀγορά, -ᾶς, ἡ"}

	for _, word := range searchResults {
		assert.Contains(t, expectedGreek, word.Greek)
	}
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
