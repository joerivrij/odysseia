//go:build !integration
// +build !integration

package app

import (
	"github.com/odysseia/plato/elastic"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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
