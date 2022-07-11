package app

import (
	"fmt"
	"github.com/kpango/glg"
	"k8s.io/api/apps/v1"
	"strconv"
	"strings"
	"time"
)

const (
	AnnotationUpdate   = "perikles/updated"
	AnnotationValidity = "perikles/validity"
	AnnotationHost     = "perikles/hostname"
	AnnotationAccesses = "perikles/accesses"
)

func (p *PeriklesHandler) checkForAnnotations(deployment v1.Deployment) error {
	for key, value := range deployment.Spec.Template.Annotations {
		if key == AnnotationHost {
			go func() {
				err := p.hostFlow(deployment)
				if err != nil {
					glg.Error(err)
				}
			}()
		}

		if key == AnnotationAccesses {
			go func() {
				err := p.clientFlow(value, deployment.Name)
				if err != nil {
					glg.Error(err)
				}
			}()
		}
	}

	return nil
}

func (p *PeriklesHandler) hostFlow(deployment v1.Deployment) error {
	var validity int
	var hostName string
	var secretName string

	for key, value := range deployment.Spec.Template.Annotations {
		glg.Info("looking through annotation")

		if key == AnnotationValidity {
			validity, _ = strconv.Atoi(value)
			glg.Info(fmt.Sprintf("validity set to %s", validity))
		}

		if key == AnnotationHost {
			hostName = value
			glg.Info(fmt.Sprintf("host set to %s", hostName))
		}

		for _, volume := range deployment.Spec.Template.Spec.Volumes {
			if volume.Secret != nil {
				secretName = volume.Secret.SecretName
			}
		}
	}

	orgName := deployment.Namespace

	hosts := []string{
		fmt.Sprintf("%s", hostName),
		fmt.Sprintf("%s.%s", hostName, orgName),
		fmt.Sprintf("%s.%s.svc", hostName, orgName),
		fmt.Sprintf("%s.%s.svc.cluster.local", hostName, orgName),
	}

	glg.Info("creating certs")
	err := p.createCert(hosts, validity, secretName)
	if err != nil {
		return err
	}

	go func() {
		err := p.syncMapping(hostName, deployment.Name, secretName, validity)
		if err != nil {
		}
	}()

	glg.Info("restarting deployment")
	err = p.restartDeployment(deployment.Namespace, deployment.Name)
	if err != nil {
		return err
	}

	return nil
}

func (p *PeriklesHandler) clientFlow(accesses, name string) error {
	hosts := strings.Split(accesses, ";")

	for _, host := range hosts {
		err := p.addClientToMapping(host, name, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PeriklesHandler) restartDeployment(ns, deploymentName string) error {
	newAnnotation := make(map[string]string)
	newAnnotation[AnnotationUpdate] = time.Now().UTC().Format(timeFormat)
	_, err := p.Config.Kube.Workload().UpdateDeploymentViaAnnotation(ns, deploymentName, newAnnotation)
	return err
}
