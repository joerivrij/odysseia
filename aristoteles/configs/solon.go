package configs

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/kubernetes"
	"github.com/odysseia-greek/plato/vault"
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
