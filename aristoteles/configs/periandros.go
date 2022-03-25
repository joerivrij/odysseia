package configs

import (
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/service"
)

type PeriandrosConfig struct {
	Namespace            string
	HttpClients          service.OdysseiaClient
	SolonCreationRequest models.SolonCreationRequest
	Kube                 kubernetes.KubeClient
}
