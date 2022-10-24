package aristoteles

import (
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/kubernetes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (c *Config) getKubeClient() (kubernetes.KubeClient, error) {
	var kubeManager kubernetes.KubeClient

	namespace := os.Getenv(EnvNamespace)
	if namespace == "" {
		namespace = defaultNamespace
	}

	kubePath := os.Getenv(EnvKubePath)

	if c.BaseConfig.OutOfClusterKube {
		var filePath string
		if kubePath == "" {
			glg.Debugf("defaulting to %s", defaultKubeConfig)
			homeDir, err := os.UserHomeDir()
			if err != nil {
				glg.Error(err)
			}

			filePath = filepath.Join(homeDir, defaultKubeConfig)
		} else {
			filePath = kubePath
		}

		cfg, err := ioutil.ReadFile(filePath)
		if err != nil {
			glg.Error("error getting kubeconfig")
		}

		kube, err := kubernetes.NewKubeClient(cfg, namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}

		kubeManager = kube
	} else {
		glg.Debug("creating in cluster kube client")
		kube, err := kubernetes.NewKubeClient(nil, namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}
		kubeManager = kube
	}

	return kubeManager, nil
}
