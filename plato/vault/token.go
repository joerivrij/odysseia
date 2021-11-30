package vault

import (
	vault "github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
)

func (v *Vault) CreateOneTimeToken(policy []string) (string, error) {
	renew := false

	tokenRequest := vault.TokenCreateRequest{
		Policies:    policy,
		TTL:         "5m",
		DisplayName: "solonCreated",
		NumUses:     1,
		Renewable:   &renew,
	}

	glg.Debug("request created")

	resp, _ := v.Connection.Auth().Token().Create(&tokenRequest)
	token := resp.Auth.ClientToken

	return token, nil
}
