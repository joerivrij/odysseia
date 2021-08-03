package elastic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthy(t *testing.T) {
	expected := "elasticsearch"
	mockCode := 200
	elasticMockClient, err := CreateMockClient("info", mockCode)
	assert.Nil(t, err)

	sut := CheckHealth(elasticMockClient)

	assert.True(t, sut.Healthy)
	assert.Equal(t, expected, sut.ClusterName)
}

func TestQueryMatchAll(t *testing.T) {
	index := "test"
	mockCode := 200
	expected := int64(3)
	elasticMockClient, err := CreateMockClient("withAll", mockCode)
	assert.Nil(t, err)

	sut, _, err := QueryWithMatchAll(*elasticMockClient, index)
	assert.Nil(t, err)
	assert.Equal(t, expected, sut.Hits.Total.Value)
}

func TestQueryMatchWithGram(t *testing.T) {
	index := "test"
	queryWord := "test"
	expected := int64(2)
	mockCode := 200
	elasticMockClient, err := CreateMockClient("withGram", mockCode)
	assert.Nil(t, err)

	sut, err := QueryMultiMatchWithGrams(*elasticMockClient, index, queryWord)
	assert.Nil(t, err)

	assert.Equal(t, expected, sut.Hits.Total.Value)
}

func TestQueryOnId(t *testing.T) {
	index := "test"
	id := "6OM7DHsBjZSpXIu1pqqS"
	mockCode := 200
	expected := int64(1)
	elasticMockClient, err := CreateMockClient("withId", mockCode)
	assert.Nil(t, err)

	sut, err := QueryOnId(*elasticMockClient, index, id)
	assert.Nil(t, err)
	assert.Equal(t, expected, sut.Hits.Total.Value)
}
