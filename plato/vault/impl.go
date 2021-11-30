package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/kpango/glg"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func GetTokenFromFile() (string, error) {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == "odysseia" {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}
	clusterKeys := filepath.Join(l, "solon", "vault_config", "cluster-keys-odysseia.json")

	f, err := ioutil.ReadFile(clusterKeys)
	if err != nil {
		glg.Error(fmt.Sprintf("Cannot read fixture file: %s", err))
		return "", err
	}

	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(f, &result)

	return result["root_token"].(string), nil
}
