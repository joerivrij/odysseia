package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

type NodesImpl struct {
	client v1.CoreV1Interface
}

func NewNodesClient(kube kubernetes.Interface) (*NodesImpl, error) {
	coreClient := kube.CoreV1()

	return &NodesImpl{client: coreClient}, nil
}

func (n *NodesImpl) List() (*corev1.NodeList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	nodes, err := n.client.Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodes, err
}
