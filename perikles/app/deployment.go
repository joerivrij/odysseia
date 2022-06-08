package app

import (
	"fmt"
	"github.com/kpango/glg"
	"time"
)

const (
	ANNOTATION_ACCESS = "perikles/secret"
	ANNOTATION_UPDATE = "perikles/updated"
)

func (p *PeriklesHandler) CheckForAnnotations(ch chan struct{}) error {
	deployments, err := p.Config.Kube.Workload().ListDeployments(p.Config.Namespace)
	if err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		for key, value := range deployment.Spec.Template.Annotations {
			if key == ANNOTATION_ACCESS {
				//this has to come from annotations?
				hostName := "perikles-service"
				orgName := deployment.Namespace
				for _, v := range deployment.Spec.Template.Spec.Volumes {
					glg.Info(v.Name)
					// get the secretname from here
					// annotation should change to server - client
					// update a secret and find all consumers
				}
				hosts := []string{
					fmt.Sprintf("%s", hostName),
					fmt.Sprintf("%s.%s", hostName, orgName),
					fmt.Sprintf("%s.%s.svc", hostName, orgName),
					fmt.Sprintf("%s.%s.svc.cluster.local", hostName, orgName),
				}

				organizations := []string{orgName}

				err := p.createCert(hosts, organizations, "perikles", value)
				if err != nil {
					ch <- struct{}{}
					return nil
				}

				err = p.restartDeployment(deployment.Namespace, deployment.Name)
				if err != nil {
					ch <- struct{}{}
					return nil
				}
			}
		}
	}

	ch <- struct{}{}

	time.Sleep(10 * time.Second)
	return nil
}

func (p *PeriklesHandler) restartDeployment(ns, deploymentName string) error {
	newAnnotation := make(map[string]string)
	newAnnotation[ANNOTATION_UPDATE] = time.Now().String()
	_, err := p.Config.Kube.Workload().UpdateDeploymentViaAnnotation(ns, deploymentName, newAnnotation)
	return err
}
