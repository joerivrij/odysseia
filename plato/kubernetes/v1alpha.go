package kubernetes

import "k8s.io/client-go/rest"

func (v *V1Alpha1Impl) ServiceMapping() ServiceMapping {
	if v == nil {
		return nil
	}
	return v.serviceMapping
}

type V1Alpha1Impl struct {
	serviceMapping *ServiceMappingsImpl
}

func NewV1AlphaClient(restConfig *rest.Config) (*V1Alpha1Impl, error) {
	serviceMapping, err := NewServiceMappingImpl(restConfig)
	if err != nil {
		return nil, err
	}

	return &V1Alpha1Impl{
		serviceMapping: serviceMapping,
	}, nil
}
