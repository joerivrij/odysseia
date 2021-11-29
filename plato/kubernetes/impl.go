package kubernetes

import (
	"context"
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

type Client interface {
	ExecNamedPod(namespace, podName string, command []string) (string, error)
	GetStatefulSets(namespace string) (*appsv1.StatefulSetList, error)
	GetPodsBySelector(namespace, selector string) (*corev1.PodList, error)
	GetPodByName(namespace, name string) (*corev1.Pod, error)
	GetDeploymentStatus(namespace string) (bool, error)
	GetSecrets(namespace string) (*corev1.SecretList, error)
	CreateSecret(namespace, secretName string, data map[string][]byte) error
	GetServiceAccounts(namespace string) (*corev1.ServiceAccountList, error)
	GetHostServer() (string, error)
	GetHostCaCert() ([]byte, error)
	CopyFileToPod(namespace, podName, destPath, srcPath string) (string, error)
}

type Kube struct {
	set    *kubernetes.Clientset
	config []byte
}

func NewKubeClient(kubeConfigFilePath string) (Client, error) {
	config, err := ioutil.ReadFile(kubeConfigFilePath)
	if err != nil {
		glg.Error("error getting kubeconfig")
		return nil, err
	}

	client, err := New(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewInClusterKubeClient(namespace string) (Client, error) {
	config, err := rest.InClusterConfig()
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	glg.Debug("created in cluster kube client")

	client := &Kube{set: clientSet, config: config.CAData }

	pods, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glg.Fatal(err.Error())
	}
	glg.Debug("There are %d pods in the cluster\n", len(pods.Items))

	return client, nil
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

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &Kube{set: clientSet, config: config }, nil
}


func (k *Kube) GetK8sClientSet() *kubernetes.Clientset {
	return k.set
}

func (k *Kube) GetConfig() []byte {
	return k.config
}

func (k *Kube)GetNewLock(lockName, podName, namespace string) *resourcelock.LeaseLock {
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