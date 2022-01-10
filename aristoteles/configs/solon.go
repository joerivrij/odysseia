package configs

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/vault"
)

type SolonConfig struct {
	Vault            vault.Client
	ElasticClient    elasticsearch.Client
	ElasticCert      []byte
	Kube             kubernetes.KubeClient
	Namespace        string
	AccessAnnotation string
	RoleAnnotation   string
}
