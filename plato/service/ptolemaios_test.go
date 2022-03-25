package service

import (
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPtolemaiosClient(t *testing.T) {
	scheme := "http"
	baseUrl := "somelocalhost.com"
	username := "testuser"
	password := "supersecretpassword"
	cert := "Somecert"
	vaultResponse := models.ElasticConfigVault{
		Username:    username,
		Password:    password,
		ElasticCERT: cert,
	}
	config := ClientConfig{
		Scheme:        scheme,
		SolonUrl:      "",
		PtolemaiosUrl: baseUrl,
	}

	t.Run("Get", func(t *testing.T) {
		codes := []int{
			200,
		}

		r, err := vaultResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		sut, err := testClient.Ptolemaios().GetSecret()
		assert.Nil(t, err)
		assert.Equal(t, password, sut.Password)
		assert.Equal(t, cert, sut.ElasticCERT)
		assert.Equal(t, username, sut.Username)
	})

}
