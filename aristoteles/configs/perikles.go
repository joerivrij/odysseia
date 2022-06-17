package configs

import (
	"github.com/odysseia/plato/certificates"
	"github.com/odysseia/plato/kubernetes"
)

type PeriklesConfig struct {
	Kube      kubernetes.KubeClient
	Cert      certificates.CertClient
	Namespace string
}
