package app

import (
	"bytes"
	"encoding/json"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/certificates"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := configs.PeriklesConfig{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/perikles/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestValidityFlow(t *testing.T) {
	ns := "test"
	organizations := []string{"test"}
	validityCa := 3650
	arJsonFilePath := filepath.Join("../fixture", "ar.json")
	jsonFile, err := os.Open(arJsonFilePath)
	assert.Nil(t, err)
	arJson, err := ioutil.ReadAll(jsonFile)
	assert.Nil(t, err)

	cert, err := certificates.NewCertGeneratorClient(organizations, validityCa)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	err = cert.InitCa()
	assert.Nil(t, err)

	t.Run("EmptyBody", func(t *testing.T) {
		fakeKube, err := kubernetes.FakeKubeClient(ns)
		assert.Nil(t, err)
		testConfig := configs.PeriklesConfig{
			Kube:      fakeKube,
			Cert:      cert,
			Namespace: ns,
		}

		expected := "request"
		router := InitRoutes(testConfig)
		response := performPostRequest(router, "/perikles/v1/validate", nil)

		var validity models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&validity)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Contains(t, validity.Messages[0].Message, expected)
	})

	t.Run("UnparseableBody", func(t *testing.T) {
		testConfig := configs.PeriklesConfig{
			Kube:      nil,
			Cert:      cert,
			Namespace: ns,
		}

		expected := "AdmissionReview"
		router := InitRoutes(testConfig)
		body := bytes.NewReader([]byte("nonvalidjson"))
		response := performPostRequest(router, "/perikles/v1/validate", body)

		var validity models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&validity)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Contains(t, validity.Messages[0].Message, expected)
	})

	t.Run("ArRequestIsNil", func(t *testing.T) {
		arNilJsonPath := filepath.Join("../fixture", "arNil.json")
		file, err := os.Open(arNilJsonPath)
		assert.Nil(t, err)
		arNilJson, err := ioutil.ReadAll(file)
		assert.Nil(t, err)
		testConfig := configs.PeriklesConfig{
			Kube:      nil,
			Cert:      cert,
			Namespace: ns,
		}

		expected := "nil request"
		router := InitRoutes(testConfig)
		bodyInBytes := bytes.NewReader(arNilJson)
		response := performPostRequest(router, "/perikles/v1/validate", bodyInBytes)

		var validity models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&validity)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Contains(t, validity.Messages[0].Message, expected)
	})

	t.Run("ValidityRequestValid", func(t *testing.T) {
		testConfig := configs.PeriklesConfig{
			Kube:      nil,
			Cert:      cert,
			Namespace: ns,
		}

		router := InitRoutes(testConfig)
		bodyInBytes := bytes.NewReader(arJson)
		response := performPostRequest(router, "/perikles/v1/validate", bodyInBytes)

		var validity v1beta1.AdmissionReview
		err = json.NewDecoder(response.Body).Decode(&validity)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, validity.Response.Allowed)
	})
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
