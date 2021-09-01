package elastic

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestUnableToCreateElasticClientWithFaultyAddress(t *testing.T) {
	os.Setenv("ELASTIC_SEARCH_SERVICE", "http://[::1]a")
	_, err := CreateElasticClientFromEnvVariables()
	if assert.Error(t, err) {
		fmt.Printf("%T\n", err)
		assert.Contains(t, err.Error(), "cannot create client:")
	}
	os.Unsetenv("ELASTIC_SEARCH_SERVICE")
}

func TestAbleToCreateElasticClientWithDefaults(t *testing.T) {
	client, err := CreateElasticClientFromEnvVariables()
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
