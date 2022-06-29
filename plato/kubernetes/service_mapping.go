package kubernetes

import (
	"context"
	"fmt"
	"github.com/odysseia/plato/kubernetes/crd/v1alpha"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	extensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

type ServiceMapping interface {
	Parse(services []v1alpha.Service, name, ns string) (*v1alpha.Mapping, error)
	List() (*v1alpha.MappingList, error)
	Get(name string) (*v1alpha.Mapping, error)
	Create(mapping *v1alpha.Mapping) (*v1alpha.Mapping, error)
	Update(mapping *v1alpha.Mapping) (*v1alpha.Mapping, error)
	CreateInCluster() (bool, error)
	GetDefinition(name string) (*apiextensionv1.CustomResourceDefinition, error)
}

type ServiceMappingsImpl struct {
	ExtensionClient *clientset.Clientset
	Client          rest.Interface
}

func NewServiceMappingImpl(kubeConfig []byte) (*ServiceMappingsImpl, error) {
	c, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}

	restConfig, err := c.ClientConfig()
	if err != nil {
		return nil, err
	}

	config := *restConfig
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha.GroupName, Version: v1alpha.Version}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	extensionClient, err := extensionsclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &ServiceMappingsImpl{
		ExtensionClient: extensionClient,
		Client:          client,
	}, nil
}

func (s *ServiceMappingsImpl) Parse(services []v1alpha.Service, name, ns string) (*v1alpha.Mapping, error) {
	apiVersion := fmt.Sprintf("%s/%s", v1alpha.GroupName, v1alpha.Version)

	if services == nil {
		services = []v1alpha.Service{}
	}

	mappingResource := v1alpha.Mapping{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		APIVersion: apiVersion,
		Kind:       v1alpha.Kind,
		Spec: v1alpha.Spec{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Services:   services,
		},
	}

	return &mappingResource, nil
}

func (s *ServiceMappingsImpl) List() (*v1alpha.MappingList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	result := v1alpha.MappingList{}
	err := s.Client.
		Get().
		Resource(v1alpha.Plural).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (s *ServiceMappingsImpl) Get(name string) (*v1alpha.Mapping, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	result := v1alpha.Mapping{}
	err := s.Client.
		Get().
		Name(name).
		Resource(v1alpha.Plural).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (s *ServiceMappingsImpl) Create(mapping *v1alpha.Mapping) (*v1alpha.Mapping, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	mappedService, err := mapping.Marshal()
	if err != nil {
		return nil, err
	}

	result := v1alpha.Mapping{}
	err = s.Client.
		Post().
		Resource(v1alpha.Plural).
		Body(mappedService).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (s *ServiceMappingsImpl) Update(mapping *v1alpha.Mapping) (*v1alpha.Mapping, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	mappedService, err := mapping.Marshal()
	if err != nil {
		return nil, err
	}

	result := v1alpha.Mapping{}
	err = s.Client.
		Put().
		Name(mapping.Name).
		Resource(v1alpha.Plural).
		Body(mappedService).
		Do(ctx).
		Into(&result)

	return &result, err
}

// CreateInCluster retrieves the Haven CRD and will deploy it to the cluster when it's not installed yet.
func (s *ServiceMappingsImpl) CreateInCluster() (bool, error) {
	var created bool
	crd := v1alpha.CreateServiceMapping()

	_, err := s.GetDefinition(crd.Name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			created = true
			_, err := s.ExtensionClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.Background(), crd, metav1.CreateOptions{})
			if err != nil {
				return created, err
			}
		}
	}

	return created, err
}

func (s *ServiceMappingsImpl) GetDefinition(name string) (*apiextensionv1.CustomResourceDefinition, error) {
	definition, err := s.ExtensionClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return definition, err
}
