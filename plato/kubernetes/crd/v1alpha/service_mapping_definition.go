package v1alpha

import (
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GroupName string = "odysseia-greek.com"
	Kind      string = "ServiceMapping"
	Version   string = "v1alpha1"
	Singular  string = "servicemapping"
	ShortName string = "servicemapping"
	Plural    string = "servicemappings"
	Name             = Plural + "." + GroupName
)

// CreateServiceMapping creates a new CRD specific to Odysseia
func CreateServiceMapping() *apiextensionv1.CustomResourceDefinition {
	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: Name},
		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: GroupName,
			Scope: apiextensionv1.ClusterScoped,
			Names: apiextensionv1.CustomResourceDefinitionNames{
				Plural:     Plural,
				Singular:   Singular,
				ShortNames: []string{ShortName},
				Kind:       Kind,
			},
			Versions: []apiextensionv1.CustomResourceDefinitionVersion{
				{
					Name:    Version,
					Served:  true,
					Storage: true,
					Schema: &apiextensionv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionv1.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]apiextensionv1.JSONSchemaProps{
										"services": {
											Type: "array",
											Items: &apiextensionv1.JSONSchemaPropsOrArray{
												Schema: &apiextensionv1.JSONSchemaProps{
													Type: "object",
													Properties: map[string]apiextensionv1.JSONSchemaProps{
														"name":       {Type: "string"},
														"namespace":  {Type: "string"},
														"kubeType":   {Type: "string"},
														"secretName": {Type: "string"},
														"active":     {Type: "boolean"},
														"created":    {Type: "string"},
														"validity":   {Type: "number"},
														"clients": {
															Type: "array",
															Items: &apiextensionv1.JSONSchemaPropsOrArray{
																Schema: &apiextensionv1.JSONSchemaProps{
																	Type: "object",
																	Properties: map[string]apiextensionv1.JSONSchemaProps{
																		"name":      {Type: "string"},
																		"namespace": {Type: "string"},
																		"kubeType":  {Type: "string"},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									Required: []string{
										"services",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return crd
}
