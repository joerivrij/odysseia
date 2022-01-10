package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
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
	SolonService *url.URL
	Kube         *kubernetes.Kube
	PodName      string
	Namespace    string
	IsPartOfJob  bool
	FullPodName  string
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

	cfgManager, _ := aristoteles.NewConfig()

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = defaultNamespace
	}

	var kubeClient *kubernetes.Kube

	if isJob {
		kube, err := cfgManager.GetKubeClient("", namespace)
		if err != nil {
			glg.Error(err)
			glg.Fatal("death has come, no kube config created")
		}

		kubeClient = kube
	}

	config := &PtolemaiosConfig{
		VaultService: vaultService,
		SolonService: solonUrl,
		Kube:         kubeClient,
		Namespace:    namespace,
		IsPartOfJob:  isJob,
		PodName:      podName,
		FullPodName:  envPodName,
	}

	return config
}
