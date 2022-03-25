package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFiles(t *testing.T) {
	expected := "odysseia"
	t.Run("CorrectPathIsReturn", func(t *testing.T) {

		sut := OdysseiaRootPath()
		assert.Contains(t, sut, expected)
	})
}
