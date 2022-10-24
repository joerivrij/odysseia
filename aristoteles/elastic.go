package aristoteles

import (
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func (c *Config) getElasticClient() (elastic.Client, error) {
	var es elastic.Client
	if c.env == "LOCAL" || c.env == "TEST" || c.BaseConfig.SidecarOverwrite {
		glg.Debug("creating local es client with tls enabled")

		esConf := c.getElasticConfig(c.BaseConfig.TLSEnabled)

		client, err := elastic.NewClient(esConf)
		if err != nil {
			glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
		}

		es = client
	} else {
		if c.BaseConfig.TLSEnabled {
			glg.Debug("getting es config from vault")
			//todo enable https for sidecars
			localHttps := false
			vaultConf, err := c.getConfigFromVault(localHttps)
			if err != nil {
				glg.Fatalf("error getting config from sidecar, shutting down: %s", err)
			}

			esConf := c.mapVaultToConf(vaultConf, c.BaseConfig.TLSEnabled)

			glg.Debug("creating es client with TLS enabled")
			client, err := elastic.NewClient(esConf)
			if err != nil {
				glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
			}

			es = client
		} else {
			esConf := c.getElasticConfig(c.BaseConfig.TLSEnabled)
			glg.Debug("creating local es client from env variables")
			client, err := elastic.NewClient(esConf)
			if err != nil {
				glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
			}

			es = client
		}
	}

	if c.BaseConfig.HealthCheck {
		standardTicks := 120 * time.Second
		tick := 1 * time.Second

		healthy := es.Health().Check(standardTicks, tick)
		if !healthy {
			glg.Fatalf("elasticClient unhealthy after %s ticks", standardTicks)
		}
	}

	return es, nil
}

func (c *Config) mapVaultToConf(vaultModel *models.ElasticConfigVault, tls bool) elastic.Config {
	elasticService := c.getElasticServiceFromEnv(tls)

	conf := elastic.Config{
		Service:     elasticService,
		Username:    vaultModel.Username,
		Password:    vaultModel.Password,
		ElasticCERT: vaultModel.ElasticCERT,
	}

	return conf
}

func (c *Config) getElasticConfig(tls bool) elastic.Config {
	elasticService := c.getElasticServiceFromEnv(tls)
	elasticUser := os.Getenv(EnvElasticUser)
	if elasticUser == "" {
		glg.Debugf("setting %s to default: %s", EnvElasticUser, elasticUsernameDefault)
		elasticUser = elasticUsernameDefault
	}
	elasticPassword := os.Getenv(EnvElasticPassword)
	if elasticPassword == "" {
		glg.Debugf("setting %s to default: %s", EnvElasticPassword, elasticPasswordDefault)
		elasticPassword = elasticPasswordDefault
	}

	var elasticCert string
	if tls {
		elasticCert = string(c.getCert())
	}

	esConf := elastic.Config{
		Service:     elasticService,
		Username:    elasticUser,
		Password:    elasticPassword,
		ElasticCERT: elasticCert,
	}

	return esConf
}

func (c *Config) getElasticServiceFromEnv(tls bool) string {
	elasticService := os.Getenv(EnvElasticService)
	if elasticService == "" {
		if tls {
			glg.Debugf("setting %s to default: %s", EnvElasticService, elasticServiceDefaultTlS)
			elasticService = elasticServiceDefaultTlS
		} else {
			glg.Debugf("setting %s to default: %s", EnvElasticService, elasticServiceDefault)
			elasticService = elasticServiceDefault
		}
	}

	return elasticService
}

func (c *Config) getCert() []byte {
	var cert []byte
	if c.env == "LOCAL" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		certPath := filepath.Join(homeDir, ".odysseia", "current", "elastic-certificate.pem")

		cert, _ = ioutil.ReadFile(certPath)

		return cert
	}

	if c.BaseConfig.TestOverwrite {
		glg.Info("trying to read cert file from file")
		path := "odysseia"
		rootPath := c.OdysseiaRootPath(path)
		certPath := filepath.Join(rootPath, "eratosthenes", "fixture", "elastic", "elastic-test-cert.pem")

		cert, _ = ioutil.ReadFile(certPath)

		return cert
	}

	glg.Info("trying to read cert file from pod")
	cert, _ = ioutil.ReadFile(certPathInPod)

	return cert
}
