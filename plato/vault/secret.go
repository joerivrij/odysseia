package vault

import (
	"fmt"
	"github.com/kpango/glg"
)

func (v *Vault)CreateNewSecret(name string, payload []byte) (bool, error) {
	vaultPath := fmt.Sprintf("configs/data/%s", name)

	secret, err := v.Connection.Logical().WriteBytes(vaultPath, payload)
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	glg.Debug(secret.Data)

	return true, nil
}
