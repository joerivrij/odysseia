package elastic

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilderClient(t *testing.T) {
	index := "test"
	term := "greek"
	searchWord := "test"

	t.Run("MatchQueryPass", func(t *testing.T) {
		file := "createDocument"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		response := testClient.Builder().MatchQuery(term, index)
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, index)
		assert.Contains(t, sut, term)
	})

	t.Run("MultiMatchWithGramPass", func(t *testing.T) {
		file := "createDocument"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		response := testClient.Builder().MultiMatchWithGram(searchWord)
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, searchWord)
	})

	t.Run("MatchAllQueryPass", func(t *testing.T) {
		file := "createDocument"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		response := testClient.Builder().MatchAll()
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, "match_all")
	})

	t.Run("CreateIndexPass", func(t *testing.T) {
		file := "createDocument"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		response := testClient.Builder().SearchAsYouTypeIndex(term)
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, term)
	})

	t.Run("CreateIndexPass", func(t *testing.T) {
		file := "createDocument"
		status := 200
		expected := "settings"
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		response := testClient.Builder().Index()
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, expected)
	})

	t.Run("CreateMultiMatch", func(t *testing.T) {
		file := "createDocument"
		status := 200
		expectedKey := "someKey"
		expectedValue := "someValue"
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		query := []map[string]string{
			{
				expectedKey: expectedValue,
			},
		}

		response := testClient.Builder().MultipleMatch(query)
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, expectedKey)
		assert.Contains(t, sut, expectedValue)
	})

	t.Run("Aggregate", func(t *testing.T) {
		file := "createDocument"
		status := 200
		field := "someField"
		aggregate := "someAggregate"
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		response := testClient.Builder().Aggregate(aggregate, field)
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, aggregate)
		assert.Contains(t, sut, field)
	})

	t.Run("FilteredAggregate", func(t *testing.T) {
		file := "createDocument"
		status := 200
		field := "someField"
		aggregate := "someAggregate"
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		response := testClient.Builder().FilteredAggregate(term, searchWord, aggregate, field)
		sut := fmt.Sprintf("%v", response)
		assert.Contains(t, sut, aggregate)
		assert.Contains(t, sut, field)
		assert.Contains(t, sut, term)
		assert.Contains(t, sut, searchWord)
	})
}
