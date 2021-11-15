package app

import (
	"github.com/kpango/glg"
	"net/url"
	"os"
)

const defaultVaultService = "http://127.0.0.1:8200"

type PtolemaiosConfig struct {
	VaultService string
	SolonService url.URL
}

func Get() *PtolemaiosConfig {
	vaultService := os.Getenv("VAULT_SERVICE")
	if vaultService == "" {
		vaultService = defaultVaultService
	}

	solonService := os.Getenv("SOLON_SERVICE")
	if solonService == "" {
		glg.Info("no connection to solon can be made")
	}

	solonUrl, _ := url.Parse(solonService)

	config := &PtolemaiosConfig{
		VaultService: vaultService,
		SolonService: *solonUrl,
	}

	return config
}
