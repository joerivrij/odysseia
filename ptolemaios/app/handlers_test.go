package app

import (
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/kubernetes"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetOneTimeToken(t *testing.T) {
	scheme := "http"
	baseUrl := "somelocalhost.com"
	token := "s.49uwenfke9fue"
	tokenResponse := models.TokenResponse{Token: token}

	config := service.ClientConfig{
		Scheme:        scheme,
		SolonUrl:      baseUrl,
		PtolemaiosUrl: "",
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

		testClient, err := service.NewFakeClient(config, codes, responses)
		assert.Nil(t, err)

		testConfig := configs.PtolemaiosConfig{
			HttpClients: testClient,
		}

		handler := PtolemaiosHandler{Config: &testConfig}

		sut, err := handler.getOneTimeToken()
		assert.Nil(t, err)
		assert.Equal(t, token, sut)
	})

	t.Run("GetWithError", func(t *testing.T) {
		codes := []int{
			500,
		}

		responses := []string{
			"error: You created",
		}

		testClient, err := service.NewFakeClient(config, codes, responses)
		assert.Nil(t, err)

		testConfig := configs.PtolemaiosConfig{
			HttpClients: testClient,
		}

		handler := PtolemaiosHandler{Config: &testConfig}

		sut, err := handler.getOneTimeToken()
		assert.NotNil(t, err)
		assert.Equal(t, "", sut)
	})
}

func TestJobExit(t *testing.T) {
	ns := "odysseia"
	expectedName := "testpod"
	duration := 10 * time.Millisecond

	t.Run("Get", func(t *testing.T) {
		testClient, err := kubernetes.FakeKubeClient(ns)
		assert.Nil(t, err)

		podSpec := kubernetes.CreatePodObjectWithExit(expectedName, ns)
		pod, err := testClient.Workload().CreatePod(ns, podSpec)
		assert.Nil(t, err)
		assert.Equal(t, pod.Name, expectedName)

		testConfig := configs.PtolemaiosConfig{
			Kube:        testClient,
			FullPodName: expectedName,
			PodName:     expectedName,
			Namespace:   ns,
		}

		handler := PtolemaiosHandler{Config: &testConfig, Duration: duration}
		jobExit := make(chan bool, 1)
		go handler.CheckForJobExit(jobExit)

		select {

		case <-jobExit:
			exitStatus := <-jobExit
			assert.True(t, exitStatus)
		}

	})
}
