package app

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/queue"
	"github.com/odysseia/aristoteles/configs"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestHandlerQueue(t *testing.T) {
	t.Run("QueueWordFound", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   "",
			Created: 0,
			Queue:   mockQueue,
		}
		duration := 10 * time.Millisecond
		testHandler := ThalesHandler{Config: &testConfig, QueueEmptyDuration: duration}
		testHandler.HandleQueue()

		assert.Equal(t, testConfig.Created, 0)
	})

	t.Run("QueueWordNotFound", func(t *testing.T) {
		file := "searchWordNoResults"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   "",
			Created: 0,
			Queue:   mockQueue,
		}
		duration := 10 * time.Millisecond
		testHandler := ThalesHandler{Config: &testConfig, QueueEmptyDuration: duration}
		testHandler.HandleQueue()

		assert.Equal(t, testConfig.Created, 0)
	})
}

func TestHandlerEmptyQueue(t *testing.T) {
	t.Run("EmptyQueue", func(t *testing.T) {
		mockQueue, err := queue.NewFakeKubeMqClient()
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: nil,
			Index:   "",
			Created: 0,
			Queue:   mockQueue,
		}

		testHandler := ThalesHandler{Config: &testConfig}
		emptyQueue := make(chan bool, 1)
		duration := 1 * time.Millisecond
		go testHandler.queueEmpty(emptyQueue, duration)
		select {

		case <-emptyQueue:
			assert.True(t, true)
		}
	})
}

func TestHandlerCreateDocuments(t *testing.T) {
	index := "test"

	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "market place",
		LinkedWord: "",
		Original:   "",
	}

	t.Run("WordIsTheSame", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		body.English = "a market place"

		testHandler := ThalesHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.True(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordWithAPronoun", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		testHandler := ThalesHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.True(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordFoundButDifferentMeaning", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		body.English = "notthesame"

		testHandler := ThalesHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordFoundDifferentMeaningWithoutAPronoun", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		body.English = "notthesame but multiple words"

		testHandler := ThalesHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		file := "searchWordNoResults"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		testHandler := ThalesHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		testHandler := ThalesHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.NotNil(t, err)
	})
}

func TestHandlerAddWord(t *testing.T) {
	index := "test"
	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "a market place",
		LinkedWord: "",
		Original:   "",
	}

	t.Run("DocumentCreated", func(t *testing.T) {
		file := "createDocument"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		testHandler := ThalesHandler{Config: &testConfig}
		testHandler.addWord(body)
		assert.Equal(t, testConfig.Created, 1)
	})

	t.Run("DocumentNotCreated", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		testHandler := ThalesHandler{Config: &testConfig}
		testHandler.addWord(body)
		assert.Equal(t, testConfig.Created, 0)
	})
}

func TestHandlerTransform(t *testing.T) {
	index := "test"
	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "a market place",
		LinkedWord: "",
		Original:   "",
	}

	t.Run("DocumentCreated", func(t *testing.T) {
		file := "createDocument"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		var wait sync.WaitGroup

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		wait.Add(1)

		testHandler := ThalesHandler{Config: &testConfig}
		testHandler.transformWord(body, &wait)
		assert.Equal(t, testConfig.Created, 1)
	})

	t.Run("DocumentNotCreated", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		var wait sync.WaitGroup

		testConfig := configs.ThalesConfig{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
			Queue:   nil,
		}

		wait.Add(1)

		testHandler := ThalesHandler{Config: &testConfig}
		testHandler.transformWord(body, &wait)
		assert.Equal(t, testConfig.Created, 0)
	})
}
