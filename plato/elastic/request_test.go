package elastic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestParsing(t *testing.T) {
	key := "key"
	value := "value"
	body := map[string]interface{}{
		key: value,
	}

	invalidBody := map[string]interface{}{
		key: make(chan int),
	}

	t.Run("Pass", func(t *testing.T) {
		response, err := toBuffer(body)
		sut := response.String()
		assert.Nil(t, err)
		assert.Contains(t, sut, key)
		assert.Contains(t, sut, value)
	})

	t.Run("Fail", func(t *testing.T) {
		response, err := toBuffer(invalidBody)
		sut := response.String()
		assert.NotNil(t, err)
		assert.Equal(t, sut, "")
	})
}
