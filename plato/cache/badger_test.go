package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBadgerClient(t *testing.T) {
	key := "testkey"
	value := "testvalue"

	t.Run("ReadValue", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		err = testClient.Set(key, value)
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.Equal(t, value, string(sut))
		testClient.Close()
	})

	t.Run("ReadEmptyValue", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, sut)

		testClient.Close()
	})
}
