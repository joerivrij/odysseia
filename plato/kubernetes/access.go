package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

type AccessImpl struct {
	serviceSet v1.ServiceAccountInterface
}

func NewAccessClient(kube kubernetes.Interface, namespace string) (*AccessImpl, error) {
	set := kube.CoreV1().ServiceAccounts(namespace)

	return &AccessImpl{serviceSet: set}, nil
}

func (a *AccessImpl) GetServiceAccounts() (*corev1.ServiceAccountList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	serviceAccounts, err := a.serviceSet.List(ctx, metav1.ListOptions{
		TypeMeta:            metav1.TypeMeta{},
		LabelSelector:       "",
		FieldSelector:       "",
		Watch:               false,
		AllowWatchBookmarks: false,
		ResourceVersion:     "",
		TimeoutSeconds:      nil,
		Limit:               0,
		Continue:            "",
	})
	if err != nil {
		return nil, err
	}

	return serviceAccounts, nil
}
