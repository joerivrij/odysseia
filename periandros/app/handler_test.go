package app

import (
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSolonHealthy(t *testing.T) {
	ns := "testNameSpace"
	scheme := "http"
	baseUrl := "somelocalhost.com"

	requestBody := models.SolonCreationRequest{
		Role:     "testrole",
		Access:   []string{"test"},
		PodName:  "somepodname",
		Username: "testuser",
	}

	healthModel := models.Health{
		Healthy:  true,
		Time:     "",
		Database: models.DatabaseHealth{},
		Memory:   models.Memory{},
	}

	config := service.ClientConfig{
		Scheme:        scheme,
		SolonUrl:      baseUrl,
		PtolemaiosUrl: "",
	}

	duration := 10 * time.Millisecond
	timeOut := 20 * time.Millisecond

	t.Run("SolonHealthy", func(t *testing.T) {

		codes := []int{
			200,
		}

		r, err := healthModel.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)

		testConfig := configs.PeriandrosConfig{
			Namespace:            ns,
			HttpClients:          testClient,
			SolonCreationRequest: requestBody,
			Kube:                 nil,
		}

		testHandler := PeriandrosHandler{Config: &testConfig, Duration: duration, Timeout: timeOut}
		healthy := testHandler.CheckSolonHealth()
		assert.True(t, healthy)
	})

	t.Run("SolonNotHealthy", func(t *testing.T) {
		codes := []int{
			500,
		}

		r, err := healthModel.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)

		testConfig := configs.PeriandrosConfig{
			Namespace:            ns,
			HttpClients:          testClient,
			SolonCreationRequest: requestBody,
			Kube:                 nil,
		}

		testHandler := PeriandrosHandler{Config: &testConfig, Duration: duration, Timeout: timeOut}
		healthy := testHandler.CheckSolonHealth()
		assert.False(t, healthy)
	})

	t.Run("SolonHealthyAfterATry", func(t *testing.T) {
		codes := []int{
			200,
			200,
		}

		notHealthy := models.Health{
			Healthy:  false,
			Time:     "",
			Database: models.DatabaseHealth{},
			Memory:   models.Memory{},
		}

		nr, err := notHealthy.Marshal()
		assert.Nil(t, err)
		r, err := healthModel.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(nr),
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)

		testConfig := configs.PeriandrosConfig{
			Namespace:            ns,
			HttpClients:          testClient,
			SolonCreationRequest: requestBody,
			Kube:                 nil,
		}

		testHandler := PeriandrosHandler{Config: &testConfig, Duration: duration, Timeout: timeOut}
		healthy := testHandler.CheckSolonHealth()
		assert.True(t, healthy)
	})
}

func TestCreatUser(t *testing.T) {
	ns := "testNameSpace"
	scheme := "http"
	baseUrl := "somelocalhost.com"

	requestBody := models.SolonCreationRequest{
		Role:     "testrole",
		Access:   []string{"test"},
		PodName:  "somepodname",
		Username: "testuser",
	}

	healthModel := models.Health{
		Healthy:  true,
		Time:     "",
		Database: models.DatabaseHealth{},
		Memory:   models.Memory{},
	}

	postResponse := models.SolonResponse{Created: true}

	config := service.ClientConfig{
		Scheme:        scheme,
		SolonUrl:      baseUrl,
		PtolemaiosUrl: "",
	}

	duration := 10 * time.Millisecond
	timeOut := 20 * time.Millisecond

	t.Run("UserCreated", func(t *testing.T) {
		codes := []int{
			200,
			201,
		}

		hr, err := healthModel.Marshal()
		assert.Nil(t, err)
		r, err := postResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(hr),
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)

		testConfig := configs.PeriandrosConfig{
			Namespace:            ns,
			HttpClients:          testClient,
			SolonCreationRequest: requestBody,
			Kube:                 nil,
		}

		testHandler := PeriandrosHandler{Config: &testConfig, Duration: duration, Timeout: timeOut}
		created, err := testHandler.CreateUser()
		assert.Nil(t, err)
		assert.True(t, created)
	})

	t.Run("UserNotCreated", func(t *testing.T) {
		codes := []int{
			200,
			500,
		}

		hr, err := healthModel.Marshal()
		assert.Nil(t, err)
		r, err := postResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(hr),
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)

		testConfig := configs.PeriandrosConfig{
			Namespace:            ns,
			HttpClients:          testClient,
			SolonCreationRequest: requestBody,
			Kube:                 nil,
		}

		testHandler := PeriandrosHandler{Config: &testConfig, Duration: duration, Timeout: timeOut}
		created, err := testHandler.CreateUser()
		assert.NotNil(t, err)
		assert.False(t, created)
	})

	t.Run("SolonNotHealthy", func(t *testing.T) {
		codes := []int{
			500,
		}

		notHealthy := models.Health{
			Healthy:  false,
			Time:     "",
			Database: models.DatabaseHealth{},
			Memory:   models.Memory{},
		}

		r, err := notHealthy.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)

		testConfig := configs.PeriandrosConfig{
			Namespace:            ns,
			HttpClients:          testClient,
			SolonCreationRequest: requestBody,
			Kube:                 nil,
		}

		testHandler := PeriandrosHandler{Config: &testConfig, Duration: duration, Timeout: timeOut}
		created, err := testHandler.CreateUser()
		assert.NotNil(t, err)
		assert.False(t, created)
	})
}
