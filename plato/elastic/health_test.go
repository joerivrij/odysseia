package elastic

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHealthClient(t *testing.T) {
	standardTicks := 10 * time.Millisecond
	tick := 10 * time.Millisecond

	t.Run("Healthy", func(t *testing.T) {
		file := "info"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		healthy := testClient.Health().Check(standardTicks, tick)
		assert.True(t, healthy)
	})

	t.Run("Unhealthy", func(t *testing.T) {
		file := "infoServiceDown"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		healthy := testClient.Health().Check(standardTicks, tick)
		assert.False(t, healthy)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		healthy := testClient.Health().Check(standardTicks, tick)

		assert.False(t, healthy)
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

		healthy := testClient.Health().Check(standardTicks, tick)

		assert.False(t, healthy)
	})
}
