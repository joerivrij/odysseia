package elastic

import (
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDocumentWithIndexClient(t *testing.T) {
	index := "test"
	body := models.Meros{
		Greek:      "μάχη",
		English:    "battle",
		LinkedWord: "",
		Original:   "",
	}

	t.Run("Created", func(t *testing.T) {
		file := "createDocument"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := body.Marshal()
		assert.Nil(t, err)

		created, err := testClient.Index().CreateDocument(index, sut)
		assert.Nil(t, err)
		assert.Equal(t, index, created.Index)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createIndex"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := body.Marshal()
		assert.Nil(t, err)

		created, err := testClient.Index().CreateDocument(index, sut)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := body.Marshal()
		assert.Nil(t, err)

		created, err := testClient.Index().CreateDocument(index, sut)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})

	t.Run("NoConnection", func(t *testing.T) {
		config := Config{
			Service:     "hhttttt://sjdsj.com",
			Username:    "",
			Password:    "",
			ElasticCERT: "",
		}
		testClient, err := NewClient(config)
		assert.Nil(t, err)

		sut, err := body.Marshal()
		assert.Nil(t, err)

		created, err := testClient.Index().CreateDocument(index, sut)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})
}

func TestCreateIndexClient(t *testing.T) {
	index := "test"
	expectedMalformed := "invalid character"
	searchWord := "someWord"

	t.Run("Created", func(t *testing.T) {
		file := "createIndex"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)

		sut, err := testClient.Index().Create(index, body)
		assert.Nil(t, err)
		assert.Equal(t, index, sut.Index)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createIndex"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)

		sut, err := testClient.Index().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)

		sut, err := testClient.Index().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), expectedMalformed)
	})

	t.Run("Unparseable", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		invalidBody := map[string]interface{}{
			"key": make(chan int),
		}

		created, err := testClient.Index().Create(index, invalidBody)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})

	t.Run("NoConnection", func(t *testing.T) {
		config := Config{
			Service:     "hhttttt://sjdsj.com",
			Username:    "",
			Password:    "",
			ElasticCERT: "",
		}
		testClient, err := NewClient(config)
		assert.Nil(t, err)

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)
		created, err := testClient.Index().Create(index, body)
		assert.NotNil(t, err)
		assert.False(t, created.Acknowledged)
	})
}

func TestDeleteIndexClient(t *testing.T) {
	index := "test"
	expectedMalformed := "invalid character"

	t.Run("Deleted", func(t *testing.T) {
		file := "deleteIndex"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "deleteIndex"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.NotNil(t, err)
		assert.False(t, sut)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.NotNil(t, err)
		assert.False(t, sut)
		assert.Contains(t, err.Error(), expectedMalformed)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.Nil(t, err)
		assert.False(t, sut)
	})

	t.Run("NoConnection", func(t *testing.T) {
		config := Config{
			Service:     "hhttttt://sjdsj.com",
			Username:    "",
			Password:    "",
			ElasticCERT: "",
		}
		testClient, err := NewClient(config)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.NotNil(t, err)
		assert.False(t, sut)
	})
}
