package app

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/vault"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultVaultService     = "http://127.0.0.1:8200"
	defaultRoleName         = "solon"
	defaultKubeConfig       = "/.kube/config"
	defaultNamespace        = "odysseia"
	defaultRoleAnnotation   = "odysseia-greek/role"
	defaultAccessAnnotation = "odysseia-greek/access"
)

type SolonConfig struct {
	Vault            vault.Client
	ElasticClient    elasticsearch.Client
	ElasticCert      []byte
	Kube             kubernetes.Kube
	Namespace        string
	AccessAnnotation string
	RoleAnnotation   string
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

	var kubeManager kubernetes.Kube
	if env != "TEST" {
		glg.Debug("creating in cluster kube client")
		kube, err := kubernetes.NewInClusterKube(namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}
		kubeManager = *kube
	} else {
		glg.Debugf("defaulting to %s", defaultKubeConfig)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			glg.Error(err)
		}

		filePath := filepath.Join(homeDir, defaultKubeConfig)

		cfg, err := ioutil.ReadFile(filePath)
		if err != nil {
			glg.Error("error getting kubeconfig")
		}

		kube, err := kubernetes.New(cfg, namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}
		kubeManager = *kube
	}

	config := &SolonConfig{
		ElasticClient:    *es,
		Vault:            vaultClient,
		ElasticCert:      cert,
		Kube:             kubeManager,
		Namespace:        namespace,
		RoleAnnotation:   defaultRoleAnnotation,
		AccessAnnotation: defaultAccessAnnotation,
	}

	return healthy, config
}
