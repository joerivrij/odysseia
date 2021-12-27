package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/configuration"
	"github.com/odysseia/plato/kubernetes"
	"net/url"
	"os"
	"strings"
)

const defaultVaultService = "http://127.0.0.1:8200"
const defaultSolonService = "http://localhost:5000"
const defaultNamespace = "odysseia"

type PtolemaiosConfig struct {
	VaultService string
	SolonService url.URL
	Kube         *kubernetes.Kube
	PodName      string
	Namespace    string
	IsPartOfJob  bool
}

func Get() *PtolemaiosConfig {
	vaultService := os.Getenv("VAULT_SERVICE")
	if vaultService == "" {
		vaultService = defaultVaultService
	}

	solonService := os.Getenv("SOLON_SERVICE")
	if solonService == "" {
		glg.Info("no connection to solon can be made")
		solonService = defaultSolonService
	}

	solonUrl, _ := url.Parse(solonService)

	envPodName := os.Getenv("POD_NAME")
	splitPodName := strings.Split(envPodName, "-")
	podName := splitPodName[0]

	var isJob bool
	job := os.Getenv("ISJOB")
	if job == "" {
		isJob = false
	} else {
		isJob = true
	}

	cfgManager, _ := configuration.NewConfig()

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = defaultNamespace
	}

	kube, err := cfgManager.GetKubeClient("", namespace)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has come, no kube config created")
	}

	config := &PtolemaiosConfig{
		VaultService: vaultService,
		SolonService: *solonUrl,
		Kube:         kube,
		Namespace:    namespace,
		IsPartOfJob:  isJob,
		PodName:      podName,
	}

	return config
}
