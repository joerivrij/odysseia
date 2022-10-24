package configs

import "github.com/odysseia-greek/plato/service"

type BaseConfig struct {
	Index                string `yaml:"INDEX"`
	SolonService         string `yaml:"SOLON_SERVICE"`
	VaultService         string `yaml:"VAULT_SERVICE"`
	TLSEnabled           bool   `yaml:"TLSENABLED"`
	Namespace            string `yaml:"NAMESPACE"`
	HealthCheckOverwrite bool   `yaml:"HEALTH_CHECK_OVERWRITE"`
	OutOfClusterKube     bool   `yaml:"OUT_OF_CLUSTER_KUBE"`
	TestOverwrite        bool   `yaml:"TEST_OVERWRITE"`
	SidecarOverwrite     bool
	HealthCheck          bool
	HttpClient           service.HttpClient
}
