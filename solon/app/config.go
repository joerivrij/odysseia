package app

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	vault "github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const defaultVaultService = "http://127.0.0.1:8200"
const defaultRoleName = "solon"

type SolonConfig struct {
	VaultClient *vault.Client
	ElasticClient elasticsearch.Client
}

func Get(ticks time.Duration, es *elasticsearch.Client) (bool, *SolonConfig) {
	healthy := elastic.CheckHealthyStatusElasticSearch(es, ticks)
	if !healthy {
		glg.Errorf("elasticClient unhealthy after %s ticks", ticks)
		return healthy, nil
	}

	var vaultClient *vault.Client
	vaultRootToken := os.Getenv("VAULT_ROOT_TOKEN")
	vaultAuthMethod := os.Getenv("AUTH_METHOD")
	vaultService := os.Getenv("VAULT_SERVICE")

	if vaultService == "" {
		vaultService = defaultVaultService
	}

	vaultRole := os.Getenv("VAULT_ROLE")
	if vaultRole == "" {
		vaultRole = defaultRoleName
	}

	glg.Debugf("vaultAuthMethod set to %s", vaultAuthMethod)

	if vaultAuthMethod == "kubernetes" {
		client, err := createVaultClientKubernetes(vaultService, vaultRole)
		if err != nil {
			glg.Error(err)
		}

		vaultClient = client
	} else {
		if vaultRootToken == "" {
			glg.Debug("root token empty getting from file for local testing")
			vaultRootToken = getTokenFromFile()
			client, err := createVaultClient(vaultService, vaultRootToken)
			if err != nil {
				glg.Error(err)
			}

			vaultClient = client
		} else {
			client, err := createVaultClient(vaultService, vaultRootToken)
			if err != nil {
				glg.Error(err)
			}
			vaultClient = client
		}
	}

	config := &SolonConfig{
		ElasticClient: *es,
		VaultClient: vaultClient,
	}

	return healthy, config
}

func createVaultClient(address, rootToken string) (*vault.Client, error) {
	config := vault.Config{
		Address:    address,
	}

	client, err := vault.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	client.SetToken(rootToken)

	return client, nil
}

func getTokenFromFile() string {
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
	clusterKeys := filepath.Join(l, "solon", "vault_config", "cluster-keys.json")

	f, err := ioutil.ReadFile(clusterKeys)
	if err != nil {
		panic(fmt.Sprintf("Cannot read fixture file: %s", err))
	}

	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(f, &result)

	return result["root_token"].(string)
}

func createVaultClientKubernetes(address, vaultRole string) (*vault.Client, error) {
	config := vault.Config{
		Address:    address,
	}

	client, err := vault.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	jwt, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, fmt.Errorf("unable to read file containing service account token: %w", err)
	}

	params := map[string]interface{}{
		"jwt":  string(jwt),
		"role": vaultRole,
	}

	// log in to Vault's Kubernetes auth method
	resp, err := client.Logical().Write("auth/kubernetes/login", params)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.ClientToken == "" {
		return nil, fmt.Errorf("login response did not return client token")
	}

	client.SetToken(resp.Auth.ClientToken)

	return client, nil
}

func InitRoot(config SolonConfig) bool {
	glg.Debug("creating secrets at startup and validating functioning of vault connection")
	roleNames := []string{"dictionary", "text", "words", "grammar"}
	indexes := []string{"dictionary", "grammar", "authors", "authorNames", "multiChoice"}
	typeOfRoles := []string{"seeder", "api"}

	var created bool
	for _, role := range roleNames {
		for _, typeOfRole := range typeOfRoles {
			glg.Debug(role)

			var privileges []string
			if typeOfRole == "seeder" {
				privileges = append(privileges, "delete")
				privileges = append(privileges, "create")
			} else {
				privileges = append(privileges, "read")
			}

			names := []string{indexes[0]}

			indices := []models.Index{
				{
					Names:      names,
					Privileges: []string{"all"},
					Query:      "",
				},
			}

			application := []models.Application{
				{
					Application: "odysseia",
					Privileges:  privileges,
					Resources:   []string{"*"},
				},
			}

			putRole := models.CreateRoleRequest{
				Cluster:      []string{"all"},
				Indices:      indices,
				Applications: application,
				RunAs:        nil,
				Metadata:     models.Metadata{Version: 1},
			}
			roleCreated, err := elastic.CreateRole(&config.ElasticClient, role, putRole)
			if err != nil {
				glg.Error(err)
			}

			created = roleCreated
		}
	}

	return created
}
