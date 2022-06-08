package configs

import "github.com/odysseia/plato/kubernetes"

type PeriklesConfig struct {
	Kube      kubernetes.KubeClient
	Namespace string
}
