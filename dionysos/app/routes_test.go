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

func TestPingPong(t *testing.T) {
	t.Run("Pass", func(t *testing.T) {
		testConfig := configs.DionysosConfig{}
		router := InitRoutes(testConfig)
		expected := "{\"result\":\"pong\"}"

		w := performGetRequest(router, "/dionysos/v1/ping")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})
}

func TestHealthEndPoint(t *testing.T) {
	t.Run("Pass", func(t *testing.T) {
		fixtureFile := "info"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic: mockElasticClient,
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/dionysos/v1/health")

		var healthModel models.Health
		err = json.NewDecoder(response.Body).Decode(&healthModel)
		assert.Nil(t, err)
		//models.Health
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, healthModel.Healthy)
	})

	t.Run("Fail", func(t *testing.T) {
		fixtureFile := "infoServiceDown"
		mockCode := 502
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic: mockElasticClient,
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/dionysos/v1/health")

		var healthModel models.Health
		err = json.NewDecoder(response.Body).Decode(&healthModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.False(t, healthModel.Healthy)
	})
}

func TestCheckGrammarEndPoint(t *testing.T) {
	t.Run("HappyPathFemFirst", func(t *testing.T) {
		fixtureFile := "dionysosFemaleNoun"
		mockCode := 200
		expected := "noun - sing - fem - nom"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, _ := QueryRuleSet(nil, "dionysos")
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic:          mockElasticClient,
			SecondaryIndex:   dictionaryIndexDefault,
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
	})

	t.Run("HappyPathMascSecond", func(t *testing.T) {
		fixtureFile := "dionysosMascNoun"
		mockCode := 200
		expected := "noun - plural - masc - nom"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, _ := QueryRuleSet(nil, "dionysos")
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic:          mockElasticClient,
			SecondaryIndex:   dictionaryIndexDefault,
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
	})

	t.Run("HappyMultiResultsFem", func(t *testing.T) {
		fixtureFile := "dionysosMultiMatch"
		mockCode := 200
		expected := "noun - sing - fem - nom"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, _ := QueryRuleSet(nil, "dionysos")
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic:          mockElasticClient,
			SecondaryIndex:   dictionaryIndexDefault,
			Index:            elasticIndexDefault,
			DeclensionConfig: *declensionConfig,
		}
		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=μάχη")

		var declensions models.DeclensionTranslationResults
		err = json.NewDecoder(response.Body).Decode(&declensions)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, len(declensions.Results) >= 1)
		assert.Equal(t, expected, declensions.Results[0].Rule)
	})

	t.Run("HappyMultiResultsMixed", func(t *testing.T) {
		fixtureFile := "dionysosMultiMatchMixed"
		mockCode := 200
		expected := "noun - sing - fem - nom"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, _ := QueryRuleSet(nil, "dionysos")
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic:          mockElasticClient,
			SecondaryIndex:   dictionaryIndexDefault,
			Index:            elasticIndexDefault,
			DeclensionConfig: *declensionConfig,
		}
		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/dionysos/v1/checkGrammar?word=μάχη")

		var declensions models.DeclensionTranslationResults
		err = json.NewDecoder(response.Body).Decode(&declensions)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, len(declensions.Results) >= 1)
		assert.Equal(t, expected, declensions.Results[0].Rule)
	})

	t.Run("HappyPathPreposition", func(t *testing.T) {
		fixtureFile := "dionysosPreposition"
		mockCode := 200
		expected := "preposition"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, _ := QueryRuleSet(nil, "dionysos")
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic:          mockElasticClient,
			SecondaryIndex:   dictionaryIndexDefault,
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
	})

	t.Run("NoQueryParam", func(t *testing.T) {
		expected := "cannot be empty"

		fixtureFile := "dionysosPreposition"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, _ := QueryRuleSet(nil, "dionysos")
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic:          mockElasticClient,
			SecondaryIndex:   dictionaryIndexDefault,
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
	})

	t.Run("NoOptionsFound", func(t *testing.T) {
		expected := "no options found"

		config := elastic.Config{
			Service:     "hhttttt://sjdsj.com",
			Username:    "",
			Password:    "",
			ElasticCERT: "",
		}
		testClient, err := elastic.NewClient(config)
		assert.Nil(t, err)

		declensionConfig, _ := QueryRuleSet(nil, "dionysos")
		assert.Nil(t, err)

		testConfig := configs.DionysosConfig{
			Elastic:          testClient,
			SecondaryIndex:   dictionaryIndexDefault,
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
	})
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
