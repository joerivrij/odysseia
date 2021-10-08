// +build !integration

package app

import (
	"github.com/odysseia/plato/elastic"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestElasticIsHealthy(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(nil, "dionysos")
	assert.Nil(t, err)

	config := Get(mockElasticClient, declensionConfig)
	assert.NotNil(t, config)
}

func TestCreateElasticConfig(t *testing.T) {
	os.Setenv("ENV", "testing")

	expected := "attic"
	fixtureFile := "declensionsDionysos"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	declensionConfig := QueryRuleSet(mockElasticClient, "dionysos")
	assert.NotNil(t, declensionConfig)

	os.Setenv("ENV", "")

	assert.Equal(t, declensionConfig.FirstDeclension.Dialect, expected)
	assert.Equal(t, declensionConfig.SecondDeclension.Dialect, expected)
}
