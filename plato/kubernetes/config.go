package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

type ConfigurationImpl struct {
	client v1.CoreV1Interface
}

func NewConfigurationClient(kube kubernetes.Interface) (*ConfigurationImpl, error) {
	coreClient := kube.CoreV1()

	return &ConfigurationImpl{client: coreClient}, nil
}

func (c *ConfigurationImpl) ListSecrets(namespace string) (*corev1.SecretList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	secrets, err := c.client.Secrets(namespace).List(ctx, metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
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

func (c *ConfigurationImpl) CreateSecret(namespace, secretName string, data map[string][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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

	_, err := c.client.Secrets(namespace).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigurationImpl) DeleteSecret(namespace, secretName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	return c.client.Secrets(namespace).Delete(ctx, secretName, metav1.DeleteOptions{})
}

func (c *ConfigurationImpl) UpdateSecret(namespace, secretName string, data map[string][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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

	_, err := c.client.Secrets(namespace).Update(ctx, &secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigurationImpl) CreateDockerSecret(namespace, secretName string, data map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
		StringData: data,
		Type:       corev1.SecretTypeDockerConfigJson,
	}

	_, err := c.client.Secrets(namespace).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigurationImpl) CreateTlSSecret(namespace, secretName string, data map[string][]byte, immutable bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Immutable: &immutable,
		Data:      data,
		Type:      corev1.SecretTypeTLS,
	}

	_, err := c.client.Secrets(namespace).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigurationImpl) UpdateTLSSecret(namespace, secretName string, data map[string][]byte, annotation map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	immutable := false

	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretName,
			Annotations: annotation,
		},
		Immutable: &immutable,
		Data:      data,
		Type:      corev1.SecretTypeTLS,
	}

	_, err := c.client.Secrets(namespace).Update(ctx, &secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GetSecret a secrets in a namespace in your kube cluster
func (c *ConfigurationImpl) GetSecret(namespace, secretName string) (*corev1.Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	secret, err := c.client.Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret, nil
}
