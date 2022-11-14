package app

import (
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/models"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHandlerCreateDocuments(t *testing.T) {
	index := "test"
	body := models.Declension{
		Name:    "firstDeclension",
		Type:    "noun",
		Dialect: "attic",
		Declensions: []models.DeclensionElement{
			{
				Declension: "-ης",
				RuleName:   "noun - sing - masc - nom",
				SearchTerm: []string{
					"ης",
				},
			},
		},
	}

	t.Run("Created", func(t *testing.T) {
		file := "createDocument"
		status := 200
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.AnaximanderConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := AnaximanderHandler{Config: &testConfig}
		err = testHandler.AddToElastic(body, &wg)
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "createIndex"
		status := 502
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.AnaximanderConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := AnaximanderHandler{Config: &testConfig}
		err = testHandler.AddToElastic(body, &wg)
		assert.NotNil(t, err)
	})
}

func TestHandlerDeleteIndex(t *testing.T) {
	index := "test"

	t.Run("Deleted", func(t *testing.T) {
		file := "deleteIndex"
		status := 201
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.AnaximanderConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := AnaximanderHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.AnaximanderConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := AnaximanderHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.NotNil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.AnaximanderConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := AnaximanderHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.NotNil(t, err)
	})
}

func TestHandlerCreateIndex(t *testing.T) {
	index := "test"

	t.Run("Created", func(t *testing.T) {
		file := "createIndex"
		status := 201
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.AnaximanderConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := AnaximanderHandler{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.AnaximanderConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := AnaximanderHandler{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.NotNil(t, err)
	})
}
