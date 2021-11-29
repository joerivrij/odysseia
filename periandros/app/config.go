package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"net/url"
	"os"
	"strings"
)

const defaultSolonService = "http://localhost:5000"

type PeriandrosConfig struct {
	SolonService url.URL
	SolonCreationRequest models.SolonCreationRequest
}

func Get() *PeriandrosConfig {
	solonService := os.Getenv("SOLON_SERVICE")
	if solonService == "" {
		glg.Info("no connection to solon can be made")
		solonService = defaultSolonService
	}

	role := os.Getenv("ELASTIC_ROLE")
	envAccess := os.Getenv("ELASTIC_ACCESS")
	podName := os.Getenv("POD_NAME")
	access := strings.Split(envAccess, ";")
	splitPodName := strings.Split(podName, "-")
	username := splitPodName[0]

	glg.Infof("username from podname is: %s", username)

	creationRequest := models.SolonCreationRequest{
		Role:    role,
		Access:  access,
		PodName: podName,
		Username: username,
	}

	solonUrl, _ := url.Parse(solonService)

	config := &PeriandrosConfig{
		SolonService: *solonUrl,
		SolonCreationRequest: creationRequest,
	}

	return config
}
