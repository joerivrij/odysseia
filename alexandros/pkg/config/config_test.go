package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidConfig(t *testing.T) {
	healthy, _ := Get(1)
	assert.False(t, healthy)
}
