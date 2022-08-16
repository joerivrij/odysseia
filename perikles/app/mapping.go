package app

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/kubernetes/crd/v1alpha"
	"time"
)

const (
	timeFormat string = "2006-01-02 15:04:05"
)

func (p *PeriklesHandler) addHostToMapping(serviceName, secretName, kubeType string, validity int) (*v1alpha.Mapping, error) {
	mapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Get(p.Config.CrdName)
	if err != nil {
		return nil, err
	}

	for i, service := range mapping.Spec.Services {
		if service.Name == serviceName {
			service.Active = true
			service.Validity = validity
			service.KubeType = kubeType
			service.SecretName = secretName
			mapping.Spec.Services[i] = service
			glg.Debugf("updating existing service mapping %s", service.Name)
			updatedMapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Update(mapping)
			if err != nil {
				return nil, err
			}

			return updatedMapping, nil
		}
	}

	service := v1alpha.Service{
		Name:       serviceName,
		KubeType:   kubeType,
		Namespace:  p.Config.Namespace,
		SecretName: secretName,
		Active:     true,
		Validity:   validity,
		Created:    time.Now().UTC().Format(timeFormat),
		Clients:    []v1alpha.Client{},
	}
	mapping.Spec.Services = append(mapping.Spec.Services, service)

	glg.Debugf("updating new service mapping %s", serviceName)
	updatedMapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Update(mapping)
	if err != nil {
		return nil, err
	}

	return updatedMapping, nil
}

func (p *PeriklesHandler) addClientToMapping(hostName, clientName, kubeType string) (*v1alpha.Mapping, error) {
	mapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Get(p.Config.CrdName)
	if err != nil {
		return nil, err
	}

	client := v1alpha.Client{
		Name:      clientName,
		KubeType:  kubeType,
		Namespace: p.Config.Namespace,
	}

	found := false
	for i, service := range mapping.Spec.Services {
		if service.Name == hostName {
			found = true
			mapping.Spec.Services[i].Clients = append(mapping.Spec.Services[i].Clients, client)
		}
	}

	if !found {
		service := v1alpha.Service{
			Name:       hostName,
			Namespace:  p.Config.Namespace,
			KubeType:   "",
			SecretName: "",
			Active:     false,
			Validity:   0,
			Created:    time.Now().UTC().Format(timeFormat),
			Clients:    []v1alpha.Client{client},
		}
		mapping.Spec.Services = append(mapping.Spec.Services, service)
	}

	updatedMapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Update(mapping)
	if err != nil {
		return nil, err
	}

	return updatedMapping, nil
}

func (p *PeriklesHandler) loopForMappingUpdates() {
	ticker := time.NewTicker(6 * time.Hour)
	for {
		select {
		case <-ticker.C:
			err := p.checkMappingForUpdates()
			if err != nil {
				glg.Error(err)
			}
		}
	}
}

func (p *PeriklesHandler) checkMappingForUpdates() error {
	mapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Get(p.Config.CrdName)
	if err != nil {
		return err
	}

	if len(mapping.Spec.Services) == 0 {
		glg.Debug("service mapping empty no action required")
		return nil
	}

	for _, service := range mapping.Spec.Services {
		redeploy, err := calculateTimeDifference(service.Validity, service.Created)
		if err != nil {
			return err
		}

		if redeploy {
			glg.Debugf("redeploy needed for service: %s", service.Name)
			glg.Debug("creating new certs after validity ran out")
			orgName := service.Namespace
			hostName := service.Name

			hosts := []string{
				fmt.Sprintf("%s", hostName),
				fmt.Sprintf("%s.%s", hostName, orgName),
				fmt.Sprintf("%s.%s.svc", hostName, orgName),
				fmt.Sprintf("%s.%s.svc.cluster.local", hostName, orgName),
			}
			err = p.createCert(hosts, service.Validity, service.SecretName)
			if err != nil {
				return err
			}

			err = p.restartKubeResource(service.Namespace, service.Name, service.KubeType)
			if err != nil {
				return err
			}

			for _, client := range service.Clients {
				go p.restartKubeResource(client.Namespace, client.Name, client.KubeType)
			}
		}
	}

	return nil
}

func calculateTimeDifference(valid int, created string) (bool, error) {
	redeploy := false
	// validity is in days recalculate to hours
	validity := valid * 24
	validFrom, err := time.Parse(timeFormat, created)
	if err != nil {
		return redeploy, err
	}

	inHours := time.Duration(validity) * time.Hour
	validTo := validFrom.Add(inHours)
	now := time.Now().UTC()

	timeDifference := validTo.Sub(now).Hours()

	if timeDifference < (time.Hour * 24).Hours() {
		redeploy = true
	}

	return redeploy, nil
}
