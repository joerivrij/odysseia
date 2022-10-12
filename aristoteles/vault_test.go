package aristoteles

import (
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/service"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetVaultConfig(t *testing.T) {
	username := "testuser"
	password := "password"
	cert := "someCert"
	body := models.ElasticConfigVault{
		Username:    username,
		Password:    password,
		ElasticCERT: cert,
	}

	baseUrl := "http://somelocalhost.com"

	t.Run("StandardConfig", func(t *testing.T) {
		os.Setenv(EnvPtolemaiosService, baseUrl)
		codes := []int{
			200,
		}

		r, err := body.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient := service.NewFakeHttpClient(responses, codes)

		cfg := Config{
			BaseConfig: configs.BaseConfig{
				HttpClient: testClient,
			},
		}

		sut, err := cfg.getConfigFromVault(false)
		assert.Nil(t, err)
		assert.Equal(t, password, sut.Password)
		assert.Equal(t, cert, sut.ElasticCERT)
		assert.Equal(t, username, sut.Username)

	})
}
