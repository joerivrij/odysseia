package kubernetes

import (
	"context"
	"github.com/kpango/glg"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

type NamespaceImpl struct {
	client v1.CoreV1Interface
}

// NewNamespaceClient to interact with Namespace interface
func NewNamespaceClient(kube kubernetes.Interface) (*NamespaceImpl, error) {
	coreClient := kube.CoreV1()

	return &NamespaceImpl{client: coreClient}, nil
}

// Create checks if a namespace exists and creates one if it does not exist yet
func (n *NamespaceImpl) Create(namespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	ns, err := n.client.Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			glg.Infof("namespace %s not found, creating", namespace)
			nsToCreate := &corev1.Namespace{
				TypeMeta: metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			ns, err = n.client.Namespaces().Create(ctx, nsToCreate, metav1.CreateOptions{})
			if err != nil {
				return err
			}

			glg.Infof("created namespace %s", ns.Name)
		}
	} else {
		glg.Infof("namespace %s already exists", ns.Name)
	}

	return err
}

// Delete removes a named namespace from your kube cluster
func (n *NamespaceImpl) Delete(namespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err := n.client.Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})

	return err
}

// List lists all namespaces within your cluster
func (n *NamespaceImpl) List() (*corev1.NamespaceList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	ns, err := n.client.Namespaces().List(ctx, metav1.ListOptions{})

	return ns, err
}
