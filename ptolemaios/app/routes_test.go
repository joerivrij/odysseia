package app

import (
	"encoding/json"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/service"
	"github.com/odysseia-greek/plato/vault"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := configs.PtolemaiosConfig{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/ptolemaios/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealth(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		mockVaultClient, err := vault.NewMockVaultClient(t)
		assert.Nil(t, err)

		testConfig := configs.PtolemaiosConfig{
			Vault: mockVaultClient,
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/ptolemaios/v1/health")

		var healthModel models.Health
		err = json.NewDecoder(response.Body).Decode(&healthModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, healthModel.Healthy)
	})
}

func TestCreateElasticSecret(t *testing.T) {
	scheme := "http"
	baseUrl := "somelocalhost.com"
	tokenResponse := models.TokenResponse{Token: ""}
	podName := "podname"
	fullPodName := "fullpodname"
	password := "iamsupersecret"
	vaultSecret := models.CreateSecretRequest{
		Data: models.ElasticConfigVault{
			Username:    podName,
			Password:    password,
			ElasticCERT: "",
		},
	}

	config := service.ClientConfig{
		Scheme:        scheme,
		SolonUrl:      baseUrl,
		PtolemaiosUrl: "",
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockVaultClient, err := vault.NewMockVaultClient(t)
		assert.Nil(t, err)
		tokenResponse.Token = mockVaultClient.GetCurrentToken()
		payload, err := vaultSecret.Marshal()
		assert.Nil(t, err)
		created, err := mockVaultClient.CreateNewSecret(podName, payload)
		assert.Nil(t, err)
		assert.True(t, created)
		codes := []int{
			200,
		}

		r, err := tokenResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)
		assert.Nil(t, err)

		testConfig := configs.PtolemaiosConfig{
			Vault:       mockVaultClient,
			HttpClients: testClient,
			PodName:     podName,
			FullPodName: fullPodName,
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/ptolemaios/v1/secret")

		var elasticModel models.ElasticConfigVault
		err = json.NewDecoder(response.Body).Decode(&elasticModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, elasticModel.Username, podName)
		assert.Equal(t, elasticModel.Password, password)
	})

	t.Run("NoTokenResponse", func(t *testing.T) {
		codes := []int{
			500,
		}

		responses := []string{
			"error creating: total error",
		}

		testClient, err := service.NewFakeClient(config, codes, responses)
		assert.Nil(t, err)

		testConfig := configs.PtolemaiosConfig{
			Vault:       nil,
			HttpClients: testClient,
			PodName:     podName,
			FullPodName: fullPodName,
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/ptolemaios/v1/secret")

		var errorModel models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&errorModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Contains(t, errorModel.Messages[0].Field, "getToken")
	})

	t.Run("VaultError", func(t *testing.T) {
		mockVaultClient, err := vault.NewMockVaultClient(t)
		assert.Nil(t, err)
		codes := []int{
			200,
		}

		r, err := tokenResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)
		assert.Nil(t, err)

		testConfig := configs.PtolemaiosConfig{
			Vault:       mockVaultClient,
			HttpClients: testClient,
			PodName:     podName,
			FullPodName: fullPodName,
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/ptolemaios/v1/secret")

		var errorModel models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&errorModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Contains(t, errorModel.Messages[0].Field, "getSecret")
	})
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
