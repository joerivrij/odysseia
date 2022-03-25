package configs

import (
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/service"
	"github.com/odysseia/plato/vault"
)

type PtolemaiosConfig struct {
	HttpClients service.OdysseiaClient
	Vault       vault.Client
	Kube        kubernetes.KubeClient
	PodName     string
	Namespace   string
	RunOnce     bool
	FullPodName string
}
