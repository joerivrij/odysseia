package configs

import (
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/vault"
)

type SolonConfig struct {
	Vault            vault.Client
	Elastic          elastic.Client
	ElasticCert      []byte
	Kube             kubernetes.KubeClient
	Namespace        string
	AccessAnnotation string
	RoleAnnotation   string
	TLSEnabled       bool
}
