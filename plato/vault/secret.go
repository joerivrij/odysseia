package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

func (v *Vault) CreateNewSecret(name string, payload []byte) (bool, error) {
	vaultPath := fmt.Sprintf("configs/data/%s", name)

	_, err := v.Connection.Logical().WriteBytes(vaultPath, payload)
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}

func (v *Vault) GetSecret(name string) (*api.Secret, error) {
	vaultPath := fmt.Sprintf("configs/data/%s", name)

	secret, err := v.Connection.Logical().Read(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read data from vault: %w", err)
	}

	return secret, nil
}
