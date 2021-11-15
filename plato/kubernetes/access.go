package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func (k *Kube) GetServiceAccounts(namespace string) (*corev1.ServiceAccountList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	serviceAccounts, err := k.GetK8sClientSet().CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{
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

