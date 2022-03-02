package configs

import (
	"github.com/odysseia/plato/kubernetes"
)

type ThrasyboulosConfig struct {
	Namespace string
	Job       string
	Kube      kubernetes.KubeClient
}
