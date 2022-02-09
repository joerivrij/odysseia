package models

type CurrentInstallConfig struct {
	ElasticPassword string `yaml:"elastic-password"`
	HarborPassword  string `yaml:"harbor-password"`
	VaultRootToken  string `yaml:"vault-root-token"`
	VaultUnsealKey  string `yaml:"vault-unseal-key"`
}

type WhiteList struct {
	AppsToInstall []string `yaml:"appsToInstall"`
}
