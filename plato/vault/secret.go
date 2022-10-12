package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
)

func (v *Vault) CreateNewSecret(name string, payload []byte) (bool, error) {
	vaultPath := fmt.Sprintf("%s/%s", v.SecretPath, name)

	_, err := v.Connection.Logical().WriteBytes(vaultPath, payload)
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}

func (v *Vault) GetSecret(name string) (*api.Secret, error) {
	vaultPath := fmt.Sprintf("%s/%s", v.SecretPath, name)
	glg.Debugf("vaultPath set to: %s", vaultPath)

	secret, err := v.Connection.Logical().Read(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read data from vault: %w", err)
	}

	return secret, nil
}
