package elastic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRoles(t *testing.T) {
	name := "test"

	t.Run("Created", func(t *testing.T) {
		file := "createRole"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		roleRequest := CreateRoleRequest{
			Cluster:      nil,
			Indices:      nil,
			Applications: nil,
			RunAs:        nil,
			Metadata:     Metadata{},
		}
		sut, err := testClient.Access().CreateRole(name, roleRequest)
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createRole"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		roleRequest := CreateRoleRequest{
			Cluster:      nil,
			Indices:      nil,
			Applications: nil,
			RunAs:        nil,
			Metadata:     Metadata{},
		}

		sut, err := testClient.Access().CreateRole(name, roleRequest)
		assert.NotNil(t, err)
		assert.False(t, sut)
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

		roleRequest := CreateRoleRequest{
			Cluster:      nil,
			Indices:      nil,
			Applications: nil,
			RunAs:        nil,
			Metadata:     Metadata{},
		}

		sut, err := testClient.Access().CreateRole(name, roleRequest)
		assert.NotNil(t, err)
		assert.False(t, sut)
	})
}

func TestCreateUser(t *testing.T) {
	name := "test"

	t.Run("Created", func(t *testing.T) {
		file := "createUser"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		userRequest := CreateUserRequest{
			Password: "password",
			Roles:    []string{"admin"},
			FullName: "Alexandros Megalos",
			Email:    "lex@megalos.com",
			Metadata: nil,
		}
		sut, err := testClient.Access().CreateUser(name, userRequest)
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createUser"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		userRequest := CreateUserRequest{
			Password: "password",
			Roles:    []string{"admin"},
			FullName: "Alexandros Megalos",
			Email:    "lex@megalos.com",
			Metadata: nil,
		}

		sut, err := testClient.Access().CreateUser(name, userRequest)
		assert.NotNil(t, err)
		assert.False(t, sut)
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

		userRequest := CreateUserRequest{
			Password: "password",
			Roles:    []string{"admin"},
			FullName: "Alexandros Megalos",
			Email:    "lex@megalos.com",
			Metadata: nil,
		}

		sut, err := testClient.Access().CreateUser(name, userRequest)
		assert.NotNil(t, err)
		assert.False(t, sut)
	})
}
