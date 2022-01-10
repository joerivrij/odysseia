package aristoteles

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/vault"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

func (c *Config) getConfigFromVault() (*models.ElasticConfigVault, error) {
	sidecarService := os.Getenv(EnvPtolemaiosService)
	if sidecarService == "" {
		glg.Info("defaulting to %s for sidecar")
		sidecarService = defaultSidecarService
	}

	u, err := url.Parse(sidecarService)
	if err != nil {
		return nil, err
	}

	u.Path = defaultSidecarPath

	response, err := helpers.GetRequest(u)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var secret models.ElasticConfigVault
	err = json.NewDecoder(response.Body).Decode(&secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func (c *Config) getVaultClient() (vault.Client, error) {
	var vaultClient vault.Client

	vaultRootToken := c.getStringFromEnv(EnvRootToken, "")
	vaultAuthMethod := c.getStringFromEnv(EnvAuthMethod, AuthMethodToken)
	vaultService := c.getStringFromEnv(EnvVaultService, c.BaseConfig.VaultService)

	vaultRole := c.getStringFromEnv(EnvVaultRole, defaultRoleName)

	glg.Debugf("vaultAuthMethod set to %s", vaultAuthMethod)

	if vaultAuthMethod == AuthMethodKube {
		jwtToken, err := os.ReadFile(serviceAccountTokenPath)
		if err != nil {
			glg.Error(err)
			return nil, err
		}

		vaultJwtToken := string(jwtToken)

		client, err := vault.CreateVaultClientKubernetes(vaultService, vaultRole, vaultJwtToken)
		if err != nil {
			glg.Error(err)
			return nil, err
		}

		vaultClient = client
	} else {
		if vaultRootToken == "" {
			glg.Debug("root token empty getting from file for local testing")
			vaultRootToken, err := c.getTokenFromFile(defaultNamespace)
			if err != nil {
				glg.Error(err)
				return nil, err
			}
			client, err := vault.CreateVaultClient(vaultService, vaultRootToken)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
		} else {
			client, err := vault.CreateVaultClient(vaultService, vaultRootToken)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
		}
	}

	if c.BaseConfig.HealthCheckOverwrite {
		healthy := vaultClient.CheckHealthyStatus(120)
		if !healthy {
			return nil, fmt.Errorf("error getting healthy status from vault")
		}
	}

	return vaultClient, nil
}

func (c *Config) getTokenFromFile(namespace string) (string, error) {
	rootPath := helpers.OdysseiaRootPath()
	clusterKeys := filepath.Join(rootPath, "solon", "vault_config", fmt.Sprintf("cluster-keys-%s.json", namespace))

	f, err := ioutil.ReadFile(clusterKeys)
	if err != nil {
		glg.Error(fmt.Sprintf("Cannot read fixture file: %s", err))
		return "", err
	}

	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(f, &result)

	return result["root_token"].(string), nil
}
