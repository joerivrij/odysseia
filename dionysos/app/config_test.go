package app

import (
	"github.com/odysseia/plato/elastic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElasticIsNotHealthy(t *testing.T) {
	esClient, _ := elastic.CreateElasticClientFromEnvVariables()
	healthy, _ := Get(1, esClient, nil)
	assert.False(t, healthy)
}

func TestElasticIsHealthy(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	healthy, _ := Get(1, mockElasticClient, declensionConfig)
	assert.True(t, healthy)
}
