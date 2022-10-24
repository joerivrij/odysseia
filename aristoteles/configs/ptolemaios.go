package configs

import (
	"github.com/odysseia-greek/plato/kubernetes"
	"github.com/odysseia-greek/plato/service"
	"github.com/odysseia-greek/plato/vault"
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
