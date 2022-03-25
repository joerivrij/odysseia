package service

import (
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSolonClient(t *testing.T) {
	scheme := "http"
	baseUrl := "somelocalhost.com"
	token := "s.49uwenfke9fue"
	tokenResponse := models.TokenResponse{Token: token}
	postResponse := models.SolonResponse{Created: true}
	config := ClientConfig{
		Scheme:        scheme,
		SolonUrl:      baseUrl,
		PtolemaiosUrl: "",
	}

	requestBody := models.SolonCreationRequest{
		Role:     "testrole",
		Access:   []string{"test"},
		PodName:  "somepodname",
		Username: "testuser",
	}

	t.Run("Get", func(t *testing.T) {
		codes := []int{
			200,
		}

		r, err := tokenResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		sut, err := testClient.Solon().OneTimeToken()
		assert.Nil(t, err)
		assert.Equal(t, token, sut.Token)
	})

	t.Run("GetWithError", func(t *testing.T) {
		codes := []int{
			500,
		}

		r, err := tokenResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		sut, err := testClient.Solon().OneTimeToken()
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("Post", func(t *testing.T) {
		codes := []int{
			201,
		}

		r, err := postResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		sut, err := testClient.Solon().Register(requestBody)
		assert.Nil(t, err)
		assert.True(t, sut.Created)
	})

	t.Run("PostWithError", func(t *testing.T) {
		codes := []int{
			500,
		}

		r, err := postResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		sut, err := testClient.Solon().Register(requestBody)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), "500")
	})
}
