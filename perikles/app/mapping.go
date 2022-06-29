package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/kubernetes/crd/v1alpha"
	"time"
)

func (p *PeriklesHandler) syncMapping(serviceName, clientName string, validity int) error {
	mapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Get(p.Config.CrdName)
	if err != nil {
		return err
	}

	addClient := true
	if clientName == "" {
		addClient = false
	}

	client := v1alpha.Client{
		Namespace: mapping.Namespace,
		Client:    clientName,
	}

	found := false
	for i, service := range mapping.Spec.Services {
		if service.Name == serviceName {
			found = true
			if addClient {
				mapping.Spec.Services[i].Clients = append(mapping.Spec.Services[i].Clients, client)
			}
		}
	}

	if !found {
		services := []v1alpha.Service{
			{
				Name:      serviceName,
				Namespace: p.Config.Namespace,
				Active:    true,
				Validity:  validity,
				Created:   time.Now().String(),
				Clients:   []v1alpha.Client{},
			},
		}

		if addClient {
			services[0].Clients = append(services[0].Clients, client)
		}
		mapping.Spec.Services = services
	}

	v1, err := p.Config.Kube.V1Alpha1().ServiceMapping().Update(mapping)
	if err != nil {
		return err
	}

	glg.Debug(v1)

	return nil
}
