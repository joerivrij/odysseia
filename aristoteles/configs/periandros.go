package configs

import (
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"net/url"
)

type PeriandrosConfig struct {
	Namespace            string
	SolonService         *url.URL
	SolonCreationRequest models.SolonCreationRequest
	Kube                 *kubernetes.KubeClient
}
