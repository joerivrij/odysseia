package vault

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecretClient(t *testing.T) {
	key := "testKey"
	value := "testValue"
	jsonBody := map[string]string{
		key: value,
	}

	name := "someTest"

	t.Run("CreateSecret", func(t *testing.T) {
		testClient, err := NewMockVaultClient(t)
		assert.Nil(t, err)

		payload, err := json.Marshal(jsonBody)
		assert.Nil(t, err)

		sut, err := testClient.CreateNewSecret(name, payload)
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("RetrieveSecret", func(t *testing.T) {
		testClient, err := NewMockVaultClient(t)
		assert.Nil(t, err)

		secret, err := testClient.GetSecret(fixtureSecretName)
		assert.Nil(t, err)
		sut := fmt.Sprintf("%v", secret.Data)
		assert.Contains(t, sut, fixtureKey)
		assert.Contains(t, sut, fixtureValue)
	})
}
