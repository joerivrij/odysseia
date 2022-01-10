package configs

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/odysseia/plato/kubernetes"
)

type DrakonConfig struct {
	Namespace     string
	PodName       string
	Kube          *kubernetes.KubeClient
	ElasticClient elasticsearch.Client
	Roles         []string
	Indexes       []string
}
