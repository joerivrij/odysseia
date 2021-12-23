package kubernetes

import (
	"github.com/kpango/glg"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type KubeClient interface {
	Access() Access
	Configuration() Configuration
	Cluster() Cluster
	Util() Util
	Workload() Workload
	Nodes() Nodes
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

type Util interface {
	CopyFileToPod(podName, destPath, srcPath string) (string, error)
}

type Workload interface {
	ExecNamedPod(namespace, podName string, command []string) (string, error)
	GetStatefulSets(namespace string) (*appsv1.StatefulSetList, error)
	GetPodsBySelector(namespace, selector string) (*corev1.PodList, error)
	GetPodByName(namespace, name string) (*corev1.Pod, error)
	GetDeploymentStatus(namespace string) (bool, error)
}

type Nodes interface {
	List() (*corev1.NodeList, error)
}

type Kube struct {
	set           *kubernetes.Clientset
	access        *AccessImpl
	util          *UtilImpl
	cluster       *ClusterImpl
	configuration *ConfigurationImpl
	workload      *WorkloadImpl
	nodes         *NodesImpl
	config        []byte
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

	cluster, err := NewClusterClient(config)
	if err != nil {
		return nil, err
	}

	configuration, err := NewConfigurationClient(clientSet)
	if err != nil {
		return nil, err
	}

	workload, err := NewWorkloadClient(clientSet)
	if err != nil {
		return nil, err
	}

	util, err := NewUtilClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	nodes, err := NewNodesClient(clientSet)
	if err != nil {
		return nil, err
	}

	return &Kube{
		set:           clientSet,
		config:        config,
		access:        access,
		cluster:       cluster,
		configuration: configuration,
		workload:      workload,
		util:          util,
		nodes:         nodes,
	}, nil
}

func NewKubeClient(filePath, ns string) (*Kube, error) {
	cfg, err := ioutil.ReadFile(filePath)
	if err != nil {
		glg.Error("error getting kubeconfig")
	}

	kube, err := New(cfg, ns)
	if err != nil {
		glg.Fatal("error creating kubeclient")
	}

	return kube, err
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

	util, err := NewUtilClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	configuration, err := NewConfigurationClient(clientSet)
	if err != nil {
		return nil, err
	}

	workload, err := NewWorkloadClient(clientSet)
	if err != nil {
		return nil, err
	}

	return &Kube{
		set:           clientSet,
		config:        config.CAData,
		access:        access,
		util:          util,
		cluster:       nil,
		configuration: configuration,
		workload:      workload,
	}, nil
}

func (k *Kube) Access() Access {
	if k == nil {
		return nil
	}
	return k.access
}

func (k *Kube) Util() Util {
	if k == nil {
		return nil
	}

	return k.util
}

func (k *Kube) Cluster() Cluster {
	if k == nil {
		return nil
	}

	return k.cluster
}

func (k *Kube) Configuration() Configuration {
	if k == nil {
		return nil
	}

	return k.configuration
}

func (k *Kube) Workload() Workload {
	if k == nil {
		return nil
	}

	return k.workload
}

func (k *Kube) Nodes() Nodes {
	if k == nil {
		return nil
	}

	return k.nodes
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
