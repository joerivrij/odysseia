package app

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateElasticConfig(t *testing.T) {
	t.Run("UsingLocal", func(t *testing.T) {
		os.Setenv("ENV", "development")

		expected := "attic"
		fixtureFile := "declensionsDionysos"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, err := QueryRuleSet(mockElasticClient, "dionysios")
		assert.NotNil(t, declensionConfig)
		assert.Nil(t, err)

		os.Setenv("ENV", "")

		assert.Equal(t, declensionConfig.Declensions[0].Dialect, expected)
		assert.Equal(t, declensionConfig.Declensions[1].Dialect, expected)
	})

	t.Run("UsingElastic", func(t *testing.T) {
		os.Setenv("ENV", "somethingelse")

		expected := "attic"
		fixtureFile := "declensionsDionysos"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, err := QueryRuleSet(mockElasticClient, "dionysios")
		assert.NotNil(t, declensionConfig)
		assert.Nil(t, err)

		os.Setenv("ENV", "")

		assert.Equal(t, declensionConfig.Declensions[0].Dialect, expected)
		assert.Equal(t, declensionConfig.Declensions[1].Dialect, expected)
	})

	t.Run("UsingElasticWithError", func(t *testing.T) {
		os.Setenv("ENV", "somethingelse")

		fixtureFile := "malformed"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		declensionConfig, err := QueryRuleSet(mockElasticClient, "dionysios")
		assert.Nil(t, declensionConfig)
		assert.NotNil(t, err)

		os.Setenv("ENV", "")
	})
}
