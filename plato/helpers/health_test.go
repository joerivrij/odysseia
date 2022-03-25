package helpers

import (
	"github.com/odysseia/plato/elastic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthClient(t *testing.T) {
	t.Run("Healthy", func(t *testing.T) {
		file := "info"
		status := 200
		testClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		sut := GetHealthOfApp(testClient)
		assert.True(t, sut.Healthy)
	})

	t.Run("HealthyWithVault", func(t *testing.T) {
		sut := GetHealthWithVault(true)
		assert.True(t, sut.Healthy)
	})
}
