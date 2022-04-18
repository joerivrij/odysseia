package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

var CliSettings = cli.New()

type HelmClient interface {
	Install(chartPath string) (*release.Release, error)
	InstallNamed(releaseName, chartPath string) (*release.Release, error)
	InstallWithValues(chartPath string, values map[string]interface{}) (*release.Release, error)
	InstallNamespaced(chartPath, namespace string, createNamespace bool) (*release.Release, error)
	List() ([]*release.Release, error)
}

type Helm struct {
	Namespace    string
	ActionConfig *action.Configuration
}

func NewHelmClient(kubeConfig []byte, ns string) (HelmClient, error) {
	c, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}

	restConfig, err := c.ClientConfig()
	if err != nil {
		return nil, err
	}

	cfg, err := kubeCliConfig(ns, restConfig)
	if err != nil {
		return nil, err
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cfg, ns, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	helm := Helm{
		Namespace:    ns,
		ActionConfig: actionConfig,
	}

	return &helm, nil
}

func kubeCliConfig(namespace string, restConfig *rest.Config) (*genericclioptions.ConfigFlags, error) {
	kubeConfig := genericclioptions.NewConfigFlags(false)
	kubeConfig.APIServer = &restConfig.Host
	kubeConfig.BearerToken = &restConfig.BearerToken
	kubeConfig.CAFile = &restConfig.CAFile
	kubeConfig.Namespace = &namespace

	return kubeConfig, nil
}
