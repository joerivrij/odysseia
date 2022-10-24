package aristoteles

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/service"
	"github.com/odysseia-greek/plato/vault"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	VAULT = "vault"
)

func (c *Config) getConfigFromVault(localHttps bool) (*models.ElasticConfigVault, error) {
	sidecarService := os.Getenv(EnvPtolemaiosService)
	if sidecarService == "" {
		glg.Infof("defaulting to %s for sidecar", defaultSidecarService)
		sidecarService = defaultSidecarService
	}

	u, err := url.Parse(sidecarService)
	if err != nil {
		return nil, err
	}

	var cert []byte
	if localHttps {
		rootPath := os.Getenv("CERT_ROOT")
		fileName := filepath.Join(rootPath, "ptolemaios")
		cert, _ = ioutil.ReadFile(fileName)
	}

	ptolemaiosClient, err := service.NewPtolemaiosConfig(u.Scheme, u.Host, cert, nil)
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
	tlsEnabled := c.getBoolFromEnv(EnvTLSEnabled)
	rootPath := c.getStringFromEnv(EnvRootTlSDir, defaultTLSFileLocation)
	secretPath := filepath.Join(rootPath, VAULT)

	glg.Debugf("vaultAuthMethod set to %s", vaultAuthMethod)
	glg.Debugf("secretPath set to %s", secretPath)
	glg.Debugf("tlsEnabled set to %v", tlsEnabled)

	var tlsConfig *api.TLSConfig

	if tlsEnabled {
		insecure := false
		if c.env == "LOCAL" || c.env == "TEST" {
			insecure = !insecure
			secretPath = "/tmp"
		}

		ca := fmt.Sprintf("%s/vault.ca", secretPath)
		cert := fmt.Sprintf("%s/vault.crt", secretPath)
		key := fmt.Sprintf("%s/vault.key", secretPath)

		tlsConfig = vault.CreateTLSConfig(insecure, ca, cert, key, secretPath)
	}

	if vaultAuthMethod == AuthMethodKube {
		jwtToken, err := os.ReadFile(serviceAccountTokenPath)
		if err != nil {
			glg.Error(err)
			return nil, err
		}

		vaultJwtToken := string(jwtToken)

		client, err := vault.CreateVaultClientKubernetes(vaultService, vaultRole, vaultJwtToken, tlsConfig)
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
			client, err := vault.NewVaultClient(vaultService, localToken, tlsConfig)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
		} else {
			client, err := vault.NewVaultClient(vaultService, vaultRootToken, tlsConfig)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
		}
	}

	glg.Debug(vaultClient)
	return vaultClient, nil
}

func (c *Config) getTokenFromFile(namespace string) (string, error) {
	path := "odysseia"
	rootPath := c.OdysseiaRootPath(path)
	clusterKeys := filepath.Join(rootPath, "eratosthenes", "fixture", "vault", fmt.Sprintf("cluster-keys-%s.json", namespace))

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

func (c *Config) OdysseiaRootPath(path string) string {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == path {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}

	return l
}
