package app

import (
	"github.com/odysseia-greek/plato/certificates"
	"github.com/odysseia-greek/plato/kubernetes"
	"github.com/odysseia/aristoteles/configs"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAnnotations(t *testing.T) {
	ns := "test"
	organizations := []string{"test"}
	validityCa := 3650
	cert, err := certificates.NewCertGeneratorClient(organizations, validityCa)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	err = cert.InitCa()
	assert.Nil(t, err)
	fakeKube, err := kubernetes.FakeKubeClient(ns)
	assert.Nil(t, err)
	crdName := "testCrd"
	testConfig := configs.PeriklesConfig{
		Kube:      fakeKube,
		Cert:      cert,
		Namespace: ns,
		CrdName:   crdName,
	}
	handler := PeriklesHandler{Config: &testConfig}
	deploymentName := "periklesDeployment"
	volumeName := "periklesVolume"
	host := "perikles"
	validity := "10"
	secretName := "periklesVolumeSecret"

	t.Run("HostOnly", func(t *testing.T) {
		annotations := map[string]string{
			AnnotationHost:     host,
			AnnotationValidity: validity,
		}
		deployment := kubernetes.CreateAnnotatedDeploymentObject(deploymentName, ns, annotations)
		err := handler.checkForAnnotations(*deployment)
		assert.Nil(t, err)
		sut, err := fakeKube.V1Alpha1().ServiceMapping().Get("asfasf")
		assert.Nil(t, err)
		found := false
		for _, service := range sut.Spec.Services {
			if service.Name == host {
				assert.Equal(t, "", service.SecretName)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("HostOnlyWithSecretFromVolume", func(t *testing.T) {
		annotations := map[string]string{
			AnnotationHost:     host,
			AnnotationValidity: validity,
		}
		deployment := kubernetes.CreateAnnotatedDeploymentObject(deploymentName, ns, annotations)
		volume := kubernetes.CreatePodSpecVolume(volumeName, secretName)
		deployment.Spec.Template.Spec.Volumes = volume
		err := handler.checkForAnnotations(*deployment)
		assert.Nil(t, err)
		sut, err := fakeKube.V1Alpha1().ServiceMapping().Get("asfasf")
		assert.Nil(t, err)
		found := false
		for _, service := range sut.Spec.Services {
			if service.Name == host {
				assert.Equal(t, secretName, service.SecretName)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("HostOnlyWithSecretFromAnnotation", func(t *testing.T) {
		hostSecret := "superSecret"
		annotations := map[string]string{
			AnnotationHost:       host,
			AnnotationValidity:   validity,
			AnnotationHostSecret: hostSecret,
		}
		deployment := kubernetes.CreateAnnotatedDeploymentObject(deploymentName, ns, annotations)
		err := handler.checkForAnnotations(*deployment)
		assert.Nil(t, err)
		sut, err := fakeKube.V1Alpha1().ServiceMapping().Get("asfasf")
		assert.Nil(t, err)
		found := false
		for _, service := range sut.Spec.Services {
			if service.Name == host {
				assert.Equal(t, hostSecret, service.SecretName)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("HostOnlyWithSecretFromAnnotation", func(t *testing.T) {
		hostSecret := "superSecret"
		annotations := map[string]string{
			AnnotationHost:       host,
			AnnotationValidity:   validity,
			AnnotationHostSecret: hostSecret,
		}
		deployment := kubernetes.CreateAnnotatedDeploymentObject(deploymentName, ns, annotations)
		err := handler.checkForAnnotations(*deployment)
		assert.Nil(t, err)
		sut, err := fakeKube.V1Alpha1().ServiceMapping().Get("asfasf")
		assert.Nil(t, err)
		found := false
		for _, service := range sut.Spec.Services {
			if service.Name == host {
				assert.Equal(t, hostSecret, service.SecretName)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("ClientOnlyWithNonExistingService", func(t *testing.T) {
		client := "archimedes;plato"
		annotations := map[string]string{
			AnnotationAccesses: client,
		}
		deployment := kubernetes.CreateAnnotatedDeploymentObject(deploymentName, ns, annotations)
		err := handler.checkForAnnotations(*deployment)
		assert.Nil(t, err)
		sut, err := fakeKube.V1Alpha1().ServiceMapping().Get("asfasf")
		assert.Nil(t, err)

		clients := strings.Split(client, ";")
		for _, service := range sut.Spec.Services {
			for _, c := range clients {
				if service.Name == c {
					assert.False(t, service.Active)
					assert.Equal(t, 1, len(service.Clients))
				}
			}
		}
	})
}
