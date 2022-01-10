package configs

import (
	"github.com/odysseia/plato/kubernetes"
	"net/url"
)

type PtolemaiosConfig struct {
	VaultService string
	SolonService *url.URL
	Kube         kubernetes.KubeClient
	PodName      string
	Namespace    string
	RunOnce      bool
	FullPodName  string
}
