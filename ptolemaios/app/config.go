package app

import (
	"github.com/kpango/glg"
	"net/url"
	"os"
)

const defaultVaultService = "http://127.0.0.1:8200"
const defaultSolonService = "http://localhost:5000"

type PtolemaiosConfig struct {
	VaultService string
	SolonService url.URL
	PodName      string
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

	podName := os.Getenv("POD_NAME")

	config := &PtolemaiosConfig{
		VaultService: vaultService,
		SolonService: *solonUrl,
		PodName:      podName,
	}

	return config
}
