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

func (p *PeriklesHandler) syncMapping(serviceName, deploymentName, secretName string, validity int) error {
	mapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Get(p.Config.CrdName)
	if err != nil {
		return err
	}

	found := false
	inactive := false
	index := 0
	for i, service := range mapping.Spec.Services {
		if service.Name == serviceName {
			if service.Active == false {
				inactive = true
				index = i
				continue
			}
			found = true
		}
	}

	if !found {
		service := v1alpha.Service{
			Name:           serviceName,
			Namespace:      p.Config.Namespace,
			DeploymentName: deploymentName,
			SecretName:     secretName,
			Active:         true,
			Validity:       validity,
			Created:        time.Now().UTC().Format(timeFormat),
			Clients:        []v1alpha.Client{},
		}

		if inactive {
			mapping.Spec.Services[index] = service
		} else {
			mapping.Spec.Services = append(mapping.Spec.Services, service)
		}
	}

	v1, err := p.Config.Kube.V1Alpha1().ServiceMapping().Update(mapping)
	if err != nil {
		return err
	}

	glg.Debug(v1)

	return nil
}

func (p *PeriklesHandler) addClientToMapping(serviceName, clientName, deploymentName string) error {
	mapping, err := p.Config.Kube.V1Alpha1().ServiceMapping().Get(p.Config.CrdName)
	if err != nil {
		return err
	}

	client := v1alpha.Client{
		Namespace:      mapping.Namespace,
		Client:         clientName,
		DeploymentName: deploymentName,
	}

	found := false
	for i, service := range mapping.Spec.Services {
		if service.Name == serviceName {
			found = true
			mapping.Spec.Services[i].Clients = append(mapping.Spec.Services[i].Clients, client)
		}
	}

	if !found {
		service := v1alpha.Service{
			Name:           serviceName,
			Namespace:      p.Config.Namespace,
			DeploymentName: deploymentName,
			SecretName:     "",
			Active:         false,
			Validity:       0,
			Created:        time.Now().UTC().Format(timeFormat),
			Clients:        []v1alpha.Client{},
		}
		service.Clients = append(service.Clients, client)
		mapping.Spec.Services = append(mapping.Spec.Services, service)
	}

	v1, err := p.Config.Kube.V1Alpha1().ServiceMapping().Update(mapping)
	if err != nil {
		return err
	}

	glg.Debug(v1)

	return nil
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

			err = p.restartDeployment(service.Namespace, service.DeploymentName)
			if err != nil {
				return err
			}

			for _, client := range service.Clients {
				go p.restartDeployment(client.Namespace, client.DeploymentName)
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
