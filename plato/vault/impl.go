package vault

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"time"
)

type Client interface {
	CheckHealthyStatus(ticks time.Duration) bool
	Health() (bool, error)
	CreateOneTimeToken(policy []string) (string, error)
	CreateNewSecret(name string, payload []byte) (bool, error)
	GetSecret(name string) (*api.Secret, error)
}

type Vault struct {
	Connection *api.Client
}

func CreateVaultClient(address, token string) (Client, error) {
	config := api.Config{
		Address: address,
	}

	client, err := api.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	client.SetToken(token)

	return &Vault{Connection: client}, nil
}

func CreateVaultClientKubernetes(address, vaultRole, jwt string) (Client, error) {
	config := api.Config{
		Address: address,
	}

	client, err := api.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	k8sAuth, err := auth.NewKubernetesAuth(
		vaultRole,
		auth.WithServiceAccountToken(jwt),
	)

	// log in to Vault's Kubernetes auth method
	resp, err := client.Auth().Login(context.TODO(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.ClientToken == "" {
		return nil, fmt.Errorf("login response did not return client token")
	}

	client.SetToken(resp.Auth.ClientToken)

	return &Vault{Connection: client}, nil
}
