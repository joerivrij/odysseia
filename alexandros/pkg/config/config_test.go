package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	healthy, _ := Get(1)
	assert.False(t, healthy)
}
