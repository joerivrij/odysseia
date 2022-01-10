package app

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/kubernetes"
	"io/ioutil"
	"os"
	"strings"
)

const (
	defaultSolonService = "http://localhost:5000"
	defaultNamespace    = "odysseia"
)

type DrakonConfig struct {
	Namespace     string
	Podname       string
	Kube          *kubernetes.Kube
	ElasticClient elasticsearch.Client
	Roles         []string
	Indexes       []string
}

func Get() *DrakonConfig {
	roles := os.Getenv("ELASTIC_ROLES")
	indexes := os.Getenv("ELASTIC_INDEXES")

	if roles == "" || indexes == "" {
		glg.Error("ELASTIC_ROLES or ELASTIC_INDEXES env variables not set!")
		glg.Fatal("cannot set access with empty env variables")
	}

	podName := os.Getenv("POD_NAME")
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = defaultNamespace
	}

	splitPodName := strings.Split(podName, "-")
	username := splitPodName[0]

	glg.Infof("username from pod is: %s", username)

	cfgManager, _ := aristoteles.NewConfig()
	kube, err := cfgManager.GetKubeClient("", namespace)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has come, no kube config created")
	}

	splitRoles := strings.Split(roles, ";")
	splitIndexes := strings.Split(indexes, ";")

	env := os.Getenv("ENV")
	if env == "" {
		env = "TEST"
	}

	var cert []byte
	var esClient *elasticsearch.Client
	if env != "TEST" {
		glg.Info("trying to read cert file from pod")
		cert, _ = ioutil.ReadFile("/app/config/certs/elastic-certificate.pem")
		es, err := elastic.CreateElasticClientFromEnvVariablesWithTLS(cert)
		if err != nil {
			glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
		}

		esClient = es
	} else {
		cert, _ = ioutil.ReadFile("/home/joerivrij/go/src/github.com/odysseia/solon/vault_config/elastic-certificate.pem")
		es, err := elastic.CreateElasticClientFromEnvVariablesWithTLS(cert)
		if err != nil {
			glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
		}

		esClient = es
	}

	config := &DrakonConfig{
		Kube:          kube,
		ElasticClient: *esClient,
		Namespace:     namespace,
		Podname:       username,
		Roles:         splitRoles,
		Indexes:       splitIndexes,
	}

	return config
}
