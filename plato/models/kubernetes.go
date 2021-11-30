package models

import "gopkg.in/yaml.v3"

func UnmarshalKubeConfig(data []byte) (KubeConfigYaml, error) {
	var y KubeConfigYaml
	err := yaml.Unmarshal(data, &y)
	return y, err
}

type KubeConfigYaml struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			CertificateAuthority     string `yaml:"certificate-authority"`
			Extensions               []struct {
				Extension struct {
					LastUpdate string `yaml:"last-update"`
					Provider   string `yaml:"provider"`
					Version    string `yaml:"version"`
				} `yaml:"extension"`
				Name string `yaml:"name"`
			} `yaml:"extensions"`
			Server string `yaml:"server"`
		} `yaml:"cluster,omitempty"`
	} `yaml:"clusters"`
	Contexts []struct {
		Name    string `yaml:"name"`
		Context struct {
			Cluster    string `yaml:"cluster"`
			Extensions []struct {
				Extension struct {
					LastUpdate string `yaml:"last-update"`
					Provider   string `yaml:"provider"`
					Version    string `yaml:"version"`
				} `yaml:"extension"`
				Name string `yaml:"name"`
			} `yaml:"extensions"`
			Namespace string `yaml:"namespace"`
			User      string `yaml:"user"`
		} `yaml:"context,omitempty"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			Exec struct {
				APIVersion         string      `yaml:"apiVersion"`
				Args               []string    `yaml:"args"`
				Command            string      `yaml:"command"`
				Env                interface{} `yaml:"env"`
				InteractiveMode    string      `yaml:"interactiveMode"`
				ProvideClusterInfo bool        `yaml:"provideClusterInfo"`
			} `yaml:"exec"`
		} `yaml:"user,omitempty"`
	} `yaml:"users"`
}
