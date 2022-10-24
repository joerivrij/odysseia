package configs

import (
	"github.com/odysseia-greek/plato/kubernetes"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/service"
)

type PeriandrosConfig struct {
	Namespace            string
	HttpClients          service.OdysseiaClient
	SolonCreationRequest models.SolonCreationRequest
	Kube                 kubernetes.KubeClient
}
