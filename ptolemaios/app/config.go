package app

import (
	"os"
)

const defaultVaultService = "http://127.0.0.1:8200"

type PtolemaiosConfig struct {
	VaultService string
}

func Get() *PtolemaiosConfig {
	vaultService := os.Getenv("VAULT_SERVICE")
	if vaultService == "" {
		vaultService = defaultVaultService
	}

	config := &PtolemaiosConfig{
		VaultService: vaultService,
	}

	return config
}
