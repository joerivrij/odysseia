package aristoteles

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/service"
	"github.com/odysseia/plato/vault"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func (c *Config) getConfigFromVault() (*models.ElasticConfigVault, error) {
	sidecarService := os.Getenv(EnvPtolemaiosService)
	if sidecarService == "" {
		glg.Infof("defaulting to %s for sidecar", defaultSidecarService)
		sidecarService = defaultSidecarService
	}

	u, err := url.Parse(sidecarService)
	if err != nil {
		return nil, err
	}

	ptolemaiosClient, err := service.NewPtolemaiosConfig(u.Scheme, u.Host, c.BaseConfig.HttpClient)
	if err != nil {
		return nil, err
	}

	glg.Debug("client created, getting secret")

	secret, err := ptolemaiosClient.GetSecret()
	if err != nil {
		return nil, err
	}

	glg.Debug("secret returned")

	return secret, nil
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

		if c.BaseConfig.HealthCheckOverwrite {
			ticks := 120 * time.Second
			tick := 1 * time.Second
			healthy := vaultClient.CheckHealthyStatus(ticks, tick)
			if !healthy {
				return nil, fmt.Errorf("error getting healthy status from vault")
			}
		}

		vaultClient = client
	} else {
		if c.env == "LOCAL" || c.env == "TEST" {
			glg.Debug("local testing, getting token from file")
			localToken, err := c.getTokenFromFile(defaultNamespace)
			if err != nil {
				glg.Error(err)
				return nil, err
			}
			client, err := vault.NewVaultClient(vaultService, localToken)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
		} else {
			client, err := vault.NewVaultClient(vaultService, vaultRootToken)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
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
