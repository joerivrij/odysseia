package kubernetes

import (
	"github.com/odysseia/plato/kubernetes/crd/v1alpha"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	defaultKubeConfig = "/.kube/config"
)

func TestCreateDefinition(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)

	filePath := filepath.Join(homeDir, defaultKubeConfig)
	cfg, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	crdClient, err := NewV1AlphaClient(cfg)
	assert.Nil(t, err)

	_, err = crdClient.ServiceMapping().CreateInCluster()
	assert.Nil(t, err)
}

func TestIntegrationDefinitionCreatedIfNotExists(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)

	filePath := filepath.Join(homeDir, defaultKubeConfig)
	cfg, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	crdClient, err := NewV1AlphaClient(cfg)
	assert.Nil(t, err)

	_, err = crdClient.ServiceMapping().CreateInCluster()
	assert.Nil(t, err)

	ns := "odysseia"
	name := "perikles"
	clients := []v1alpha.Client{
		{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Namespace:  ns,
			Client:     "ptolemaios",
		},
		{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Namespace:  ns,
			Client:     "anotherservice",
		},
	}

	service := []v1alpha.Service{
		{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Name:       "solon",
			Namespace:  ns,
			Active:     true,
			Created:    time.Now().String(),
			Clients:    clients,
		},
	}
	mapping, err := crdClient.ServiceMapping().Parse(service, name, ns)
	assert.Nil(t, err)

	result, err := crdClient.ServiceMapping().Create(mapping)
	assert.Nil(t, err)
	assert.Equal(t, name, result.Name)
}
