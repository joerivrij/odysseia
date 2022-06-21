package kubernetes

import (
	"github.com/odysseia/plato/models"
	"io/ioutil"
)

type ClusterImpl struct {
	kubeConfig []byte
}

func NewClusterClient(config []byte) (*ClusterImpl, error) {
	return &ClusterImpl{kubeConfig: config}, nil
}

func (c *ClusterImpl) GetHostServer() (string, error) {
	var server string
	var configCluster string

	config, err := models.UnmarshalKubeConfig(c.kubeConfig)
	if err != nil {
		return "", err
	}

	currentCtx := config.CurrentContext

	for _, context := range config.Contexts {
		if context.Name == currentCtx {
			configCluster = context.Context.Cluster
			break
		}
	}

	for _, cluster := range config.Clusters {
		if cluster.Name == configCluster {
			server = cluster.Cluster.Server
		}
	}

	return server, nil
}

func (c *ClusterImpl) GetHostCaCert() ([]byte, error) {
	config, err := models.UnmarshalKubeConfig(c.kubeConfig)
	if err != nil {
		return nil, err
	}

	currentCtx := config.CurrentContext
	var configCluster string
	for _, context := range config.Contexts {
		if context.Name == currentCtx {
			configCluster = context.Context.Cluster
			break
		}
	}

	for _, cluster := range config.Clusters {
		if cluster.Name == configCluster {
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
