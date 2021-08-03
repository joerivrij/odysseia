package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomizer(t *testing.T) {
	topNumber := 10
	sut := false
	randomNumber := GenerateRandomNumber(topNumber)

	if randomNumber <= topNumber {
		sut = true
	}

	assert.True(t, sut)
}
