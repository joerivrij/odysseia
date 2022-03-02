package kubernetes

import (
	"github.com/kpango/glg"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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
	Namespaces() Namespace
}

type Access interface {
	GetServiceAccounts() (*corev1.ServiceAccountList, error)
}

type Configuration interface {
	GetSecret(namespace, secretName string) (*corev1.Secret, error)
	ListSecrets(namespace string) (*corev1.SecretList, error)
	CreateSecret(namespace, secretName string, data map[string][]byte) error
	CreateDockerSecret(namespace, secretName string, data map[string]string) error
}

type Cluster interface {
	GetHostServer() (string, error)
	GetHostCaCert() ([]byte, error)
}

type Namespace interface {
	Create(namespace string) error
	Delete(namespace string) error
	List() (*corev1.NamespaceList, error)
}

type Util interface {
	CopyFileToPod(podName, destPath, srcPath string) (string, error)
	CopyFileFromPod(srcPath, destPath, namespace, podName string) error
}

type Workload interface {
	List(namespace string) (*corev1.PodList, error)
	CreatePodSpec(namespace, name, podImage string, command []string) *corev1.Pod
	DeletePod(namespace, podName string) error
	CreatePod(namespace string, pod *corev1.Pod) (*corev1.Pod, error)
	ExecNamedPod(namespace, podName string, command []string) (string, error)
	GetStatefulSets(namespace string) (*appsv1.StatefulSetList, error)
	GetPodsBySelector(namespace, selector string) (*corev1.PodList, error)
	GetPodByName(namespace, name string) (*corev1.Pod, error)
	GetDeploymentStatus(namespace string) (bool, error)
	GetJob(namespace, name string) (*batchv1.Job, error)
	ListJobs(namespace string) (*batchv1.JobList, error)
	GetNewLock(lockName, podName, namespace string) *resourcelock.LeaseLock
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
	namespace     *NamespaceImpl
	config        []byte
}

func NewKubeClient(cfg []byte, ns string) (KubeClient, error) {
	var kube *Kube
	var err error

	inCluster := true

	if cfg != nil {
		inCluster = false
	}

	if inCluster {
		kube, err = NewInClusterKube(ns)
		if err != nil {
			return nil, err
		}
	} else {
		kube, err = NewConfigBasedKube(cfg, ns)
		if err != nil {
			return nil, err
		}
	}

	return kube, err
}

func NewConfigBasedKube(config []byte, ns string) (*Kube, error) {
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

	namespaceClient, err := NewNamespaceClient(clientSet)
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
		namespace:     namespaceClient,
	}, nil
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

func (k *Kube) Namespaces() Namespace {
	if k == nil {
		return nil
	}

	return k.namespace
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

//func (k *Kube) GetK8sClientSet() *kubernetes.Clientset {
//	return k.set
//}
//
//func (k *Kube) GetConfig() []byte {
//	return k.config
//}
