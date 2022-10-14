package configs

import (
	"github.com/odysseia-greek/plato/kubernetes"
)

type ThrasyboulosConfig struct {
	Namespace string
	Job       string
	Kube      kubernetes.KubeClient
}
