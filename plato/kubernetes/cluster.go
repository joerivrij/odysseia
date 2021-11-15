package kubernetes

import (
	"github.com/odysseia/plato/models"
	"io/ioutil"
)

func (k *Kube) GetHostServer() (string, error) {
	var server string
	kubeConfig := k.GetConfig()
	config, err := models.UnmarshalKubeConfig(kubeConfig)
	if err != nil {
		return "", err
	}

	currentCtx := config.CurrentContext

	for _, cluster := range config.Clusters {
		if cluster.Name == currentCtx {
			server = cluster.Cluster.Server
		}
	}

	return server, nil
}

func (k *Kube) GetHostCaCert() ([]byte, error) {
	kubeConfig := k.GetConfig()
	config, err := models.UnmarshalKubeConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	currentCtx := config.CurrentContext

	for _, cluster := range config.Clusters {
		if cluster.Name == currentCtx {
			if cluster.Cluster.CertificateAuthorityData == "" {
				filePath := cluster.Cluster.CertificateAuthority
				content, err := ioutil.ReadFile(filePath)
				if err != nil {
					return nil, err
				}
				return content, nil
			} else {
				return []byte(cluster.Cluster.CertificateAuthorityData), nil
			}
		}
	}

	return nil, nil
}

