package app

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/vault"
	"os"
	"path/filepath"
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
	Vault vault.Client
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

	var vaultClient vault.Client
	vaultRootToken := os.Getenv("VAULT_ROOT_TOKEN")
	vaultAuthMethod := os.Getenv("AUTH_METHOD")
	vaultService := os.Getenv("VAULT_SERVICE")
	vaultJwtToken := os.Getenv("VAULT_JWT")

	if vaultService == "" {
		vaultService = defaultVaultService
	}

	vaultRole := os.Getenv("VAULT_ROLE")
	if vaultRole == "" {
		vaultRole = defaultRoleName
	}

	glg.Debugf("vaultAuthMethod set to %s", vaultAuthMethod)

	if vaultAuthMethod == "kubernetes" {
		client, err := vault.CreateVaultClientKubernetes(vaultService, vaultRole, vaultJwtToken)
		if err != nil {
			glg.Error(err)
		}

		healthy := client.CheckHealthyStatus(120)
		if !healthy {
			glg.Fatal("death has found me")
		}
		vaultClient = client
	} else {
		if vaultRootToken == "" {
			glg.Debug("root token empty getting from file for local testing")
			vaultRootToken, err := vault.GetTokenFromFile()
			if err != nil {
				glg.Error(err)
			}
			client, err := vault.CreateVaultClient(vaultService, vaultRootToken)
			if err != nil {
				glg.Error(err)
			}

			healthy := client.CheckHealthyStatus(120)
			if !healthy {
				glg.Fatal("death has found me")
			}
			vaultClient = client
		} else {
			client, err := vault.CreateVaultClient(vaultService, vaultRootToken)
			if err != nil {
				glg.Error(err)
			}

			healthy := client.CheckHealthyStatus(120)
			if !healthy {
				glg.Fatal("death has found me")
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
		Vault: vaultClient,
		ElasticCert: cert,
		Kube: kubeManager,
		Namespace: namespace,
		RoleAnnotation: defaultRoleAnnotation,
		AccessAnnotation: defaultAccessAnnotation,
	}

	return healthy, config
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
