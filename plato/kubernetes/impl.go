package kubernetes

import (
	"github.com/kpango/glg"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type KubeClient interface {
	Access
	Configuration
	Cluster
	Copier
	Workload
}

type Access interface {
	GetServiceAccounts() (*corev1.ServiceAccountList, error)
}

type Configuration interface {
	GetSecrets(namespace string) (*corev1.SecretList, error)
	CreateSecret(namespace, secretName string, data map[string][]byte) error
}

type Cluster interface {
	GetHostServer() (string, error)
	GetHostCaCert() ([]byte, error)
}

type Copier interface {
	CopyFileToPod(podName, destPath, srcPath string) (string, error)
}

type Workload interface {
	ExecNamedPod(namespace, podName string, command []string) (string, error)
	GetStatefulSets(namespace string) (*appsv1.StatefulSetList, error)
	GetPodsBySelector(namespace, selector string) (*corev1.PodList, error)
	GetPodByName(namespace, name string) (*corev1.Pod, error)
	GetDeploymentStatus(namespace string) (bool, error)
}

type Kube struct {
	set    *kubernetes.Clientset
	Access Access
	Copier Copier
	config []byte
}

func New(config []byte, ns string) (*Kube, error) {
	c, err := clientcmd.NewClientConfigFromBytes(config)
	if err != nil {
		return nil, err
	}

	restConfig, err := c.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	access, err := NewAccessClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	return &Kube{set: clientSet, config: config, Access: access}, nil
}

func NewInClusterKube(ns string) (*Kube, error) {
	config, err := rest.InClusterConfig()
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	glg.Debug("created in cluster kube client")

	access, err := NewAccessClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	copier, err := NewCopierClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	return &Kube{set: clientSet, config: config.CAData, Access: access, Copier: copier}, nil
}

func (k *Kube) GetK8sClientSet() *kubernetes.Clientset {
	return k.set
}

func (k *Kube) GetConfig() []byte {
	return k.config
}

func (k *Kube) GetNewLock(lockName, podName, namespace string) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockName,
			Namespace: namespace,
		},
		Client: k.GetK8sClientSet().CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podName,
		},
	}
}
