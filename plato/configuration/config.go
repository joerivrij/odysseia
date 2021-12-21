package configuration

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Config interface {
	GetElasticClient() (*elasticsearch.Client, error)
	GetKubeClient(kubePath, namespace string) (*kubernetes.Kube, error)
}

const (
	defaultSidecarService = "http://127.0.0.1:5001"
	defaultSidecarPath    = "/ptolemaios/v1/secret"
	defaultKubeConfig     = "/.kube/config"
	defaultNamespace      = "odysseia"
)

type ConfigImpl struct {
	env        string
	tlsEnabled bool
	sideCar    url.URL
}

func NewConfig() (Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "TEST"
	}

	var tls bool
	envTls := os.Getenv("TLSENABLED")
	if envTls == "" {
		tls = false
	} else if envTls == "true" || envTls == "yes" {
		tls = true
	} else {
		tls = false
	}

	return &ConfigImpl{env: env,
		tlsEnabled: tls}, nil
}

func (c *ConfigImpl) GetElasticClient() (*elasticsearch.Client, error) {
	var es *elasticsearch.Client
	if c.tlsEnabled {
		esConf, err := c.GetSecretFromVault()
		if err != nil {
			glg.Fatalf("error getting config from sidecar, shutting down: %s", err)
		}

		client, err := elastic.CreateElasticClientWithTlS(*esConf)
		if err != nil {
			glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
		}

		es = client
	}

	if c.env == "TEST" {
		client, err := elastic.CreateElasticClientFromEnvVariables()
		if err != nil {
			glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
		}

		es = client
	}

	standardTicks := time.Minute * 2

	healthy := elastic.CheckHealthyStatusElasticSearch(es, standardTicks)
	if !healthy {
		glg.Fatalf("elasticClient unhealthy after %s ticks", standardTicks)
	}

	return es, nil
}

func (c *ConfigImpl) GetSecretFromVault() (*models.ElasticConfigVault, error) {
	sidecarService := os.Getenv("PTOLEMAIOS_SERVICE")
	if sidecarService == "" {
		glg.Info("defaulting to %s for sidecar")
		sidecarService = defaultSidecarService
	}

	u, err := url.Parse(sidecarService)
	if err != nil {
		return nil, err
	}

	u.Path = defaultSidecarPath

	response, err := helpers.GetRequest(*u)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var secret models.ElasticConfigVault
	err = json.NewDecoder(response.Body).Decode(&secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func (c *ConfigImpl) GetKubeClient(kubePath, namespace string) (*kubernetes.Kube, error) {
	var kubeManager kubernetes.Kube

	if namespace == "" {
		namespace = defaultNamespace
	}

	if c.env == "ARCHIMEDES" || c.env == "TEST" {
		var filePath string
		if kubePath == "" {
			glg.Debugf("defaulting to %s", defaultKubeConfig)
			homeDir, err := os.UserHomeDir()
			if err != nil {
				glg.Error(err)
			}

			filePath = filepath.Join(homeDir, defaultKubeConfig)
		} else {
			filePath = kubePath
		}

		cfg, err := ioutil.ReadFile(filePath)
		if err != nil {
			glg.Error("error getting kubeconfig")
		}

		kube, err := kubernetes.New(cfg, namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}

		kubeManager = *kube
	} else {
		glg.Debug("creating in cluster kube client")
		kube, err := kubernetes.NewInClusterKube(namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}
		kubeManager = *kube
	}

	return &kubeManager, nil
}
