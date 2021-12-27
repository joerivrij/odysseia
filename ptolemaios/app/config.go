package app

import (
	"github.com/kpango/glg"
	"net/url"
	"os"
	"strings"
)

const defaultVaultService = "http://127.0.0.1:8200"
const defaultSolonService = "http://localhost:5000"

type PtolemaiosConfig struct {
	VaultService string
	SolonService url.URL
	PodName      string
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

	config := &PtolemaiosConfig{
		VaultService: vaultService,
		SolonService: *solonUrl,
		IsPartOfJob:  isJob,
		PodName:      podName,
	}

	return config
}
