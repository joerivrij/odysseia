package app

import (
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/elastic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerCreateDocuments(t *testing.T) {
	t.Run("CreateRole", func(t *testing.T) {
		file := "createRole"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		roles := []string{seederRole, hybridRole, apiRole}

		for _, role := range roles {
			testConfig := configs.DrakonConfig{
				Elastic: mockElasticClient,
				Indexes: []string{"test"},
				Roles:   []string{role},
			}

			testHandler := DrakonHandler{Config: &testConfig}
			created, err := testHandler.CreateRoles()
			assert.Nil(t, err)
			assert.True(t, created)
		}
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createRole"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.DrakonConfig{
			Elastic: mockElasticClient,
			Indexes: []string{"test"},
			Roles:   []string{"rike"},
		}

		testHandler := DrakonHandler{Config: &testConfig}
		created, err := testHandler.CreateRoles()
		assert.NotNil(t, err)
		assert.False(t, created)
	})
}
