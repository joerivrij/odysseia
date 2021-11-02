package kubernetes

import (
	"context"
	"github.com/kpango/glg"
	"io/ioutil"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)
type KubeClient interface {
	ExecNamedPod(namespace, podName string, command []string) (string, error)
	GetStatefulSets(namespace string) (*v12.StatefulSetList, error)
	GetPodsBySelector(namespace, selector string) (*v1.PodList, error)
	GetDeploymentStatus(namespace string) (bool, error)
	GetSecrets(namespace string) (*v1.SecretList, error)
	CreateSecret(namespace, secretName string, data map[string][]byte) error
	GetServiceAccounts(namespace string) (*v1.ServiceAccountList, error)
	GetHostServer() (string, error)
	GetHostCaCert() ([]byte, error)
	CopyFileToPod(namespace, podName, destPath, srcPath string) (string, error)
}

type Kube struct {
	clientset *kubernetes.Clientset
	config    []byte
	ctx       context.Context
}

func NewKubeClient(kubeConfigFilePath string) (KubeClient, error) {
	config, err := ioutil.ReadFile(kubeConfigFilePath)
	if err != nil {
		glg.Error("error getting kubeconfig")
		return nil, err
	}

	client, err := New(config)
	if err != nil {
		return nil, err
	}
	return &KubeClientImpl{kubeClient: client}, nil
}

func New(config []byte) (*Kube, error) {
	c, err := clientcmd.NewClientConfigFromBytes(config)
	if err != nil {
		return nil, err
	}

	restConfig, err := c.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &Kube{clientset: clientset, config: config, ctx: context.Background()}, nil
}

func (k *Kube) GetK8sClientSet() *kubernetes.Clientset {
	return k.clientset
}

func (k *Kube) GetConfig() []byte {
	return k.config
}


