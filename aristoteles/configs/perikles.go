package configs

import (
	"github.com/odysseia-greek/plato/certificates"
	"github.com/odysseia-greek/plato/kubernetes"
)

type PeriklesConfig struct {
	Kube      kubernetes.KubeClient
	Cert      certificates.CertClient
	Namespace string
	CrdName   string
	TLSFiles  string
}
