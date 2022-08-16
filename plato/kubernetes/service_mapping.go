package kubernetes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia/plato/kubernetes/crd/v1alpha"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	extensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	fakeclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	manualfake "k8s.io/client-go/rest/fake"
	"net/http"
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
	ExtensionClient clientset.Interface
	Client          rest.Interface
}

const (
	timeFormat string = "2006-01-02 15:04:05"
)

var updatedMapping v1alpha.Mapping

func NewServiceMappingImpl(restConfig *rest.Config) (*ServiceMappingsImpl, error) {
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

func NewFakeServiceMappingImpl() (*ServiceMappingsImpl, error) {
	ns := "odysseia"
	mappingName := "testCrd"
	kubeType := "Deployment"
	clientName := "client"
	serviceName := "fakedService"
	validity := 10

	header := http.Header{}
	header.Set("Content-Type", runtime.ContentTypeJSON)

	fakeClient := &manualfake.RESTClient{
		GroupVersion:         appsv1.SchemeGroupVersion,
		NegotiatedSerializer: scheme.Codecs,
		Client: manualfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			mapping := v1alpha.Mapping{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: mappingName,
				},
				APIVersion: fmt.Sprintf("%s/%s", v1alpha.GroupName, v1alpha.Version),
				Kind:       v1alpha.Kind,
				Spec: v1alpha.Spec{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Services: []v1alpha.Service{
						{
							TypeMeta:   metav1.TypeMeta{},
							ObjectMeta: metav1.ObjectMeta{},
							Name:       serviceName,
							KubeType:   kubeType,
							SecretName: "",
							Namespace:  ns,
							Active:     true,
							Created:    time.Now().UTC().Format(timeFormat),
							Validity:   validity,
							Clients: []v1alpha.Client{
								{
									TypeMeta:   metav1.TypeMeta{},
									ObjectMeta: metav1.ObjectMeta{},
									Namespace:  ns,
									Name:       clientName,
									KubeType:   kubeType,
								},
							},
						},
					},
				},
			}

			switch req.Method {
			case "GET":
				byteResponse, _ := mapping.Marshal()
				if updatedMapping.APIVersion != "" {
					byteResponse, _ = updatedMapping.Marshal()
				}
				return &http.Response{StatusCode: http.StatusOK, Header: header, Body: ioutil.NopCloser(bytes.NewReader(byteResponse))}, nil
			case "PUT":
				var v1map v1alpha.Mapping
				err := json.NewDecoder(req.Body).Decode(&v1map)
				if err == nil {
					updatedMapping = v1map
				}
				byteResponse, _ := updatedMapping.Marshal()
				return &http.Response{StatusCode: http.StatusOK, Header: header, Body: ioutil.NopCloser(bytes.NewReader(byteResponse))}, nil
			case "POST":
				return &http.Response{StatusCode: http.StatusOK, Header: header, Body: ioutil.NopCloser(req.Body)}, nil
			default:
				fmt.Errorf("unexpected request: %#v\n%#v", req.URL, req)
				return nil, nil
			}
		}),
	}
	clientConfig := &rest.Config{
		APIPath: "/apis",
		ContentConfig: rest.ContentConfig{
			NegotiatedSerializer: scheme.Codecs,
			GroupVersion:         &appsv1.SchemeGroupVersion,
		},
	}
	restClient, _ := rest.RESTClientFor(clientConfig)
	restClient.Client = fakeClient.Client

	extensionClient := fakeclientset.NewSimpleClientset()

	return &ServiceMappingsImpl{
		ExtensionClient: extensionClient,
		Client:          restClient,
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

// CreateInCluster retrieves the CRD and will deploy it to the cluster when it's not installed yet.
func (s *ServiceMappingsImpl) CreateInCluster() (bool, error) {
	var created bool
	crd := v1alpha.CreateServiceMapping()

	_, err := s.GetDefinition(crd.Name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			created = true
			_, err = s.ExtensionClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.Background(), crd, metav1.CreateOptions{})
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
