package app

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/queue"
	"github.com/odysseia/aristoteles/configs"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHandlerCreateDocuments(t *testing.T) {
	index := "test"
	body := models.Logos{Logos: []models.Word{{
		Method:      "",
		Category:    "",
		Greek:       "ἀγγέλλω",
		Translation: "to bear a message ",
		Chapter:     0,
	},
	},
	}

	queueItem := true
	method := "testmethod"
	category := "testcategory"

	t.Run("CreatedWithQueue", func(t *testing.T) {
		file := "createDocument"
		status := 200
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.Add(body, &wg, method, category, queueItem)
		assert.Nil(t, err)
	})

	t.Run("CreatedWithoutQueue", func(t *testing.T) {
		file := "createDocument"
		status := 200
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		addToQueue := false
		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.Add(body, &wg, method, category, addToQueue)
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "createIndex"
		status := 502
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.Add(body, &wg, method, category, queueItem)
		assert.NotNil(t, err)
	})
}

func TestHandlerQueue(t *testing.T) {
	t.Run("QueueReturnsAnError", func(t *testing.T) {
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: nil,
			Index:   "",
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.Queue(nil)
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
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
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
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ParmenidesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.NotNil(t, err)
	})
}
