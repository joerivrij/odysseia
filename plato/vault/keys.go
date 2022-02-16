package vault

import (
	"github.com/odysseia/plato/models"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetClusterKeys() (*models.CurrentInstallConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fileName := "config.yaml"
	clusterKeyFilePath := filepath.Join(homeDir, ".odysseia", "current", fileName)
	f, err := ioutil.ReadFile(clusterKeyFilePath)
	if err != nil {
		return nil, err
	}

	var currentKeys models.CurrentInstallConfig
	err = yaml.Unmarshal(f, &currentKeys)
	if err != nil {
		return nil, err
	}

	return &currentKeys, nil
}
