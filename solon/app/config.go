package app

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	vault "github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	defaultVaultService = "http://127.0.0.1:8200"
	defaultRoleName = "solon"
	defaultKubeConfig = "/.kube/config"
	defaultNamespace = "odysseia"
	defaultRoleAnnotation = "odysseia-greek/role"
	defaultAccessAnnotation = "odysseia-greek/access"
)

type SolonConfig struct {
	VaultClient *vault.Client
	ElasticClient elasticsearch.Client
	ElasticCert []byte
	Kube kubernetes.Client
	Namespace string
	AccessAnnotation string
	RoleAnnotation string
}

func Get(ticks time.Duration, es *elasticsearch.Client, cert []byte, env string) (bool, *SolonConfig) {
	healthy := elastic.CheckHealthyStatusElasticSearch(es, ticks)
	if !healthy {
		glg.Errorf("elasticClient unhealthy after %s ticks", ticks)
		return healthy, nil
	}

	namespace := os.Getenv("NAMESPACE")

	if namespace == "" {
		namespace = defaultNamespace
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

	var kubeManager kubernetes.Client
	if env != "TEST" {
		glg.Debug("creating in cluster kube client")
		kube, err := kubernetes.NewInClusterKubeClient()
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}
		kubeManager = kube
	} else {
		glg.Debugf("defaulting to %s", defaultKubeConfig)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			glg.Error(err)
		}

		filePath := filepath.Join(homeDir, defaultKubeConfig)
		kube, err := kubernetes.NewKubeClient(filePath)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}
		kubeManager = kube
	}

	config := &SolonConfig{
		ElasticClient: *es,
		VaultClient: vaultClient,
		ElasticCert: cert,
		Kube: kubeManager,
		Namespace: namespace,
		RoleAnnotation: defaultRoleAnnotation,
		AccessAnnotation: defaultAccessAnnotation,
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
	clusterKeys := filepath.Join(l, "solon", "vault_config", "cluster-keys-odysseia.json")

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

	unparsedIndexes := os.Getenv("ELASTIC_INDEXES")
	unparsedRoles := os.Getenv("ELASTIC_ROLES")
	roles := strings.Split(unparsedRoles, ";")
	indexes := strings.Split(unparsedIndexes, ";")


	var created bool
	for _, index := range indexes {
		for _, role := range roles {
			glg.Debugf("creating a role for index %s with role %s", index, role)

			var privileges []string
			if role == "seeder" {
				privileges = append(privileges, "delete")
				privileges = append(privileges, "create")
			} else {
				privileges = append(privileges, "read")
			}

			names := []string{index}

			indices := []models.Index{
				{
					Names:      names,
					Privileges: privileges,
					Query:      "",
				},
			}

			application := []models.Application{
			}

			putRole := models.CreateRoleRequest{
				Cluster:      []string{"all"},
				Indices:      indices,
				Applications: application,
				RunAs:        nil,
				Metadata:     models.Metadata{Version: 1},
			}

			roleName := fmt.Sprintf("%s_%s", index, role)
			roleCreated, err := elastic.CreateRole(&config.ElasticClient, roleName, putRole)
			if err != nil {
				glg.Error(err)
			}

			created = roleCreated
		}
	}

	return created
}
