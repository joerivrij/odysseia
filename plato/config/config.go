package config

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/models"
	"net/url"
	"os"
	"strings"
)

type Config interface {
	GetConfigFromSidecar() (*models.SecretData, error)
}

const defaultSidecarService = "http://localhost:5001"
const defaultSidecarPath = "/ptolemaios/v1/secret"

type Sidecar struct {
	SecretName string
	Sidecar url.URL
}

func NewConfBuilderWithSidecar() (Config, error) {
	sidecarService := os.Getenv("PTOLEMAIOS_SERVICE")
	if sidecarService == "" {
		glg.Info("defaulting to %s for sidecar")
		sidecarService = defaultSidecarService
	}

	u, err := url.Parse(sidecarService)
	if err != nil {
		return nil, err
	}
	u.Path = defaultSidecarPath

	podName := os.Getenv("POD_NAME")
	splitPodName := strings.Split(podName, "-")
	secretName := splitPodName[0]

	return &Sidecar{
		SecretName: secretName,
		Sidecar: *u,
	}, nil
}

func (c *Sidecar)GetConfigFromSidecar() (*models.SecretData, error){
	response, err := helpers.GetRequest(c.Sidecar)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var secret models.SecretData
	err = json.NewDecoder(response.Body).Decode(&secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}
