package app

import (
	"fmt"
	"github.com/kpango/glg"
	"k8s.io/api/apps/v1"
	"strconv"
	"time"
)

const (
	AnnotationAccess   = "perikles/access"
	AnnotationUpdate   = "perikles/updated"
	AnnotationValidity = "perikles/validity"
	AnnotationHost     = "perikles/hostname"
	SERVER             = "server"
)

func (p *PeriklesHandler) checkForAnnotations(deployment v1.Deployment) error {
	for key, value := range deployment.Spec.Template.Annotations {
		if key == AnnotationAccess {
			glg.Info("looking through annotation")
			var validity int
			var hostName string
			var secretName string
			for k, v := range deployment.Spec.Template.Annotations {
				if k == AnnotationValidity {
					validity, _ = strconv.Atoi(v)
					glg.Info(fmt.Sprintf("validity set to %s", v))
				}

				if k == AnnotationHost {
					hostName = v
					glg.Info(fmt.Sprintf("host set to %s", hostName))
				}
			}

			orgName := deployment.Namespace
			for _, volume := range deployment.Spec.Template.Spec.Volumes {
				if volume.Secret != nil {
					secretName = volume.Secret.SecretName
				}
			}
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

			clientName := ""
			if value != SERVER {
				clientName = value
			}

			go func() {
				err := p.syncMapping(hostName, clientName, validity)
				if err != nil {
				}
			}()

			glg.Info("restarting deployment")
			err = p.restartDeployment(deployment.Namespace, deployment.Name)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (p *PeriklesHandler) restartDeployment(ns, deploymentName string) error {
	newAnnotation := make(map[string]string)
	newAnnotation[AnnotationUpdate] = time.Now().String()
	_, err := p.Config.Kube.Workload().UpdateDeploymentViaAnnotation(ns, deploymentName, newAnnotation)
	return err
}
