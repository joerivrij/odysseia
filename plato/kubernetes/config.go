package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func (k *Kube) GetSecrets(namespace string) (*corev1.SecretList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	secrets, err := k.GetK8sClientSet().CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
		TypeMeta:            metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
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

	return secrets, nil
}

func (k *Kube) CreateSecret(namespace, secretName string, data map[string][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Immutable:  nil,
		Data:       data,
		StringData: nil,
		Type:       corev1.SecretTypeOpaque,
	}

	_, err := k.GetK8sClientSet().CoreV1().Secrets(namespace).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
