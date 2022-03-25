package configs

import (
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/kubernetes"
)

type DrakonConfig struct {
	Namespace string
	PodName   string
	Kube      kubernetes.KubeClient
	Elastic   elastic.Client
	Roles     []string
	Indexes   []string
}
