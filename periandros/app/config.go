package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/configuration"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"net/url"
	"os"
	"strings"
)

const (
	defaultSolonService = "http://localhost:5000"
	defaultNamespace    = "odysseia"
)

type PeriandrosConfig struct {
	Namespace            string
	SolonService         url.URL
	SolonCreationRequest models.SolonCreationRequest
	Kube                 *kubernetes.Kube
}

func Get() *PeriandrosConfig {
	solonService := os.Getenv("SOLON_SERVICE")
	if solonService == "" {
		glg.Infof("no solon service select defaulting to %s", defaultSolonService)
		solonService = defaultSolonService
	}

	role := os.Getenv("ELASTIC_ROLE")
	envAccess := os.Getenv("ELASTIC_ACCESS")

	if role == "" || envAccess == "" {
		glg.Error("ELASTIC_ROLE or ELASTIC_ACCESS env variables not set!")
		glg.Fatal("cannot set access with empty env variables")
	}
	podName := os.Getenv("POD_NAME")
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = defaultNamespace
	}
	access := strings.Split(envAccess, ";")
	splitPodName := strings.Split(podName, "-")
	username := splitPodName[0]

	glg.Infof("username from pod is: %s", username)

	creationRequest := models.SolonCreationRequest{
		Role:     role,
		Access:   access,
		PodName:  podName,
		Username: username,
	}

	solonUrl, _ := url.Parse(solonService)

	cfgManager, _ := configuration.NewConfig()
	kube, err := cfgManager.GetKubeClient("", namespace)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has come, no kube config created")
	}

	config := &PeriandrosConfig{
		Kube:                 kube,
		Namespace:            namespace,
		SolonService:         *solonUrl,
		SolonCreationRequest: creationRequest,
	}

	return config
}
