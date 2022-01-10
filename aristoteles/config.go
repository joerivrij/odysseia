package aristoteles

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/vault"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type Config struct {
	env        string
	BaseConfig configs.BaseConfig
}

//go:embed base
var base embed.FS

func NewConfig(v interface{}) interface{} {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ARISTOTELES
	glg.Info("\n  ____  ____   ____ _____ ______   ___   ______    ___  _        ___  _____\n /    ||    \\ |    / ___/|      | /   \\ |      |  /  _]| |      /  _]/ ___/\n|  o  ||  D  ) |  (   \\_ |      ||     ||      | /  [_ | |     /  [_(   \\_ \n|     ||    /  |  |\\__  ||_|  |_||  O  ||_|  |_||    _]| |___ |    _]\\__  |\n|  _  ||    \\  |  |/  \\ |  |  |  |     |  |  |  |   [_ |     ||   [_ /  \\ |\n|  |  ||  .  \\ |  |\\    |  |  |  |     |  |  |  |     ||     ||     |\\    |\n|__|__||__|\\_||____|\\___|  |__|   \\___/   |__|  |_____||_____||_____| \\___|\n                                                                           \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"Τριών δει παιδεία: φύσεως, μαθήσεως, ασκήσεως.\"")
	glg.Info("\"Education needs these three: natural endowment, study, practice.\"")
	glg.Info(strings.Repeat("~", 37))

	envDir := os.Getenv(EnvKey)
	if envDir == "" {
		envDir = "local"
	}

	glg.Infof("getting config: %s from yaml files", envDir)

	config, err := base.ReadFile(fmt.Sprintf("%s/%s/%s", baseDir, envDir, configFileName))
	if err != nil {
		glg.Error(err)
		glg.Fatal("could not read base config")
	}

	var baseConfig configs.BaseConfig
	err = yaml.Unmarshal(config, &baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("could not unmarshal base config")
	}

	envTls := os.Getenv(EnvTlSKey)
	if envTls != "" {
		glg.Infof("%s overwritten by env variable", EnvTlSKey)
		if envTls == "true" || envTls == "yes" {
			baseConfig.TLSEnabled = true
		} else {
			baseConfig.TLSEnabled = false
		}
	}

	newConfig := Config{BaseConfig: baseConfig}
	newConfig.env = strings.ToUpper(envDir)

	glg.Infof("env set to: %s", newConfig.env)

	healthCheck := true
	healthCheckOverwrite := os.Getenv(EnvHealthCheckOverwrite)
	if healthCheckOverwrite == "yes" || healthCheckOverwrite == "true" {
		healthCheck = false
	}

	if newConfig.BaseConfig.HealthCheckOverwrite {
		healthCheck = !newConfig.BaseConfig.HealthCheckOverwrite
	}

	newConfig.BaseConfig.HealthCheck = healthCheck

	glg.Info("validating config")
	valid, err := newConfig.configValidator(v)
	if !valid || err != nil {
		glg.Error(err)
		glg.Fatal("invalid config found, shutting down")
	}

	glg.Info("config is valid")
	glg.Info("setting remaining fields")
	readConfig, err := newConfig.New(v)
	if err != nil {
		glg.Error(err)
		glg.Fatal("could not create config, shutting down")
	}

	glg.Info("config created returning")

	return readConfig
}

func (c *Config) New(v interface{}) (interface{}, error) {
	cfg := reflect.New(reflect.TypeOf(v)).Interface()
	elements := reflect.ValueOf(cfg).Elem()
	c.fillFields(&elements)

	return cfg, nil
}

func (c *Config) configValidator(v interface{}) (bool, error) {
	elements := reflect.TypeOf(v)
	for i := 0; i < elements.NumField(); i++ {
		fieldName := elements.Field(i).Name
		inList := false
		for _, validField := range validFields {
			if validField == fieldName {
				inList = true
				break
			}
		}

		if !inList {
			return false, fmt.Errorf("value: %s could not be found in valid fields, configuration file malformed", fieldName)
		}

	}

	return true, nil
}

func (c *Config) fillFields(e *reflect.Value) {
	for i := 0; i < e.NumField(); i++ {
		fieldName := e.Type().Field(i).Name
		fieldType := e.Type().Field(i).Type

		if fieldType.Kind() == reflect.String {
			switch fieldName {
			case "Index":
				indexName := c.getStringFromEnv(EnvIndex, c.BaseConfig.Index)
				e.FieldByName(fieldName).SetString(indexName)
			case "PodName":
				podName := c.getParsedPodNameFromEnv()
				e.FieldByName(fieldName).SetString(podName)
			case "FullPodName":
				podName := c.getStringFromEnv(EnvPodName, defaultPodName)
				e.FieldByName(fieldName).SetString(podName)
			case "Namespace":
				ns := c.getStringFromEnv(EnvNamespace, c.BaseConfig.Namespace)
				e.FieldByName(fieldName).SetString(ns)
			case "VaultService":
				vs := c.getStringFromEnv(EnvVaultService, c.BaseConfig.VaultService)
				e.FieldByName(fieldName).SetString(vs)
			case "SearchWord":
				vs := c.getStringFromEnv(EnvSearchWord, defaultSearchWord)
				e.FieldByName(fieldName).SetString(vs)
			case "DictionaryIndex":
				e.FieldByName(fieldName).SetString(defaultDictionaryIndex)
			case "RoleAnnotation":
				e.FieldByName(fieldName).SetString(defaultRoleAnnotation)
			case "AccessAnnotation":
				e.FieldByName(fieldName).SetString(defaultAccessAnnotation)
			}
		}

		if fieldType.Kind() == reflect.Slice {
			switch fieldName {
			case "Roles":
				roles := c.getSliceFromEnv(EnvRoles)
				rRoles := reflect.ValueOf(roles)
				e.FieldByName(fieldName).Set(rRoles)

			case "Indexes":
				indexes := c.getSliceFromEnv(EnvIndexes)
				rIndexes := reflect.ValueOf(indexes)
				e.FieldByName(fieldName).Set(rIndexes)

			case "ElasticCert":
				cert := c.getCert()
				rCert := reflect.ValueOf(cert)
				e.FieldByName(fieldName).Set(rCert)
			}
		}

		if fieldType.Kind() == reflect.Bool {
			switch fieldName {
			case "RunOnce":
				vb := c.getBoolFromEnv(EnvRunOnce)
				e.FieldByName(fieldName).SetBool(vb)
			}
		}

		switch fieldType {
		case reflect.TypeOf(elasticsearch.Client{}):

			es, err := c.getElasticClient(c.BaseConfig.HealthCheck)
			if err != nil {
				glg.Fatal("error getting es config")
				esv := reflect.ValueOf(es)
				e.FieldByName(fieldName).Set(esv)
			}

		case reflect.TypeOf((*kubernetes.KubeClient)(nil)):
			k, err := c.getKubeClient("", "")
			if err != nil {
				glg.Fatal("error getting kubeconfig")
			}
			kv := reflect.ValueOf(k)
			e.FieldByName(fieldName).Set(kv)

		case reflect.TypeOf((*vault.Client)(nil)):
			vault, err := c.getVaultClient()
			if err != nil {
				glg.Fatal("error getting vaultClient")
			}
			vv := reflect.ValueOf(vault)
			e.FieldByName(fieldName).Set(vv)

		case reflect.TypeOf((*url.URL)(nil)):
			var defaultValue string
			elem := reflect.ValueOf(&c.BaseConfig).Elem()
			for i := 0; i < elem.NumField(); i++ {
				innerFieldName := elem.Type().Field(i).Name
				if innerFieldName == fieldName {
					defaultValue = elem.Field(i).String()
					break
				}
			}

			u, _ := c.getUrl(fieldName, defaultValue)
			uv := reflect.ValueOf(u)
			e.FieldByName(fieldName).Set(uv)

		case reflect.TypeOf(models.SolonCreationRequest{}):
			request := c.getInitCreation()
			rv := reflect.ValueOf(request)
			e.FieldByName(fieldName).Set(rv)
		}
	}

}

func (c *Config) getParsedPodNameFromEnv() string {
	envPodName := os.Getenv(EnvPodName)
	if envPodName == "" {
		glg.Debugf("%s empty set as env variable - defaulting to %s", EnvPodName, defaultPodName)
		envPodName = defaultPodName
	}
	splitPodName := strings.Split(envPodName, "-")
	podName := splitPodName[0]

	return podName
}

func (c *Config) getStringFromEnv(envName, defaultValue string) string {
	var value string
	value = os.Getenv(envName)
	if value == "" {
		glg.Debugf("%s empty set as env variable - defaulting to %s", envName, defaultValue)
		value = defaultValue
	}

	return value
}

func (c *Config) getSliceFromEnv(sliceName string) []string {
	slice := os.Getenv(sliceName)

	if slice == "" {
		glg.Error("ELASTIC_ROLES or ELASTIC_INDEXES env variables not set!")
	}

	splitSlice := strings.Split(slice, ";")

	return splitSlice
}

func (c *Config) getBoolFromEnv(envName string) bool {
	var value bool
	envValue := os.Getenv(envName)
	if envValue == "" {
		value = false
	} else {
		value = true
	}

	return value
}

func (c *Config) getUrl(serviceName, defaultValue string) (*url.URL, error) {
	var envVar string
	for key, value := range serviceMapping {
		if key == serviceName {
			envVar = value
		}
	}
	service := os.Getenv(envVar)
	if service == "" {
		service = defaultValue
	}

	u, err := url.Parse(service)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *Config) getInitCreation() models.SolonCreationRequest {
	role := c.getStringFromEnv(EnvRole, "")
	envAccess := c.getSliceFromEnv(EnvIndex)
	podName := c.getStringFromEnv(EnvPodName, defaultPodName)
	splitPodName := strings.Split(podName, "-")
	username := splitPodName[0]

	glg.Infof("username from pod is: %s", username)

	creationRequest := models.SolonCreationRequest{
		Role:     role,
		Access:   envAccess,
		PodName:  podName,
		Username: username,
	}

	return creationRequest
}

func (c *Config) getElasticClient(healthCheck bool) (*elasticsearch.Client, error) {
	var es *elasticsearch.Client
	if c.env == "LOCAL" {
		if c.BaseConfig.TLSEnabled {
			glg.Debug("creating local es client with tls enabled")

			elasticCert := c.getCert()
			esConf := c.getElasticConfig(c.BaseConfig.TLSEnabled, elasticCert)

			client, err := elastic.CreateElasticClientWithTlS(esConf)
			if err != nil {
				glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
			}

			es = client
		} else {
			esConf := c.getElasticConfig(c.BaseConfig.TLSEnabled, nil)
			glg.Debug("creating local es client from env variables")
			client, err := elastic.CreateElasticClientFromEnvVariables(esConf)
			if err != nil {
				glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
			}

			es = client
		}
	} else {
		if c.BaseConfig.TLSEnabled {
			glg.Debug("getting es config from vault")
			esConf, err := c.getConfigFromVault()
			if err != nil {
				glg.Fatalf("error getting config from sidecar, shutting down: %s", err)
			}

			glg.Debug("creating es client with TLS enabled")
			client, err := elastic.CreateElasticClientWithTlS(*esConf)
			if err != nil {
				glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
			}

			es = client
		} else {
			glg.Debug("creating local es client from env variables")
			client, err := elastic.CreateElasticClientFromEnvVariables()
			if err != nil {
				glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
			}

			es = client
		}
	}

	if healthCheck {
		standardTicks := 120 * time.Second

		healthy := elastic.CheckHealthyStatusElasticSearch(es, standardTicks)
		if !healthy {
			glg.Fatalf("elasticClient unhealthy after %s ticks", standardTicks)
		}
	}

	return es, nil
}

func (c *Config) getElasticConfig(tls bool, cert []byte) models.ElasticConfig {
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
	if cert != nil {
		elasticCert = string(cert)
	} else {
		elasticCert = ""
	}

	esConf := models.ElasticConfig{
		Service:     elasticService,
		Username:    elasticUser,
		Password:    elasticPassword,
		ElasticCERT: elasticCert,
	}

	return esConf
}

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

func (c *Config) getKubeClient(kubePath, namespace string) (*kubernetes.KubeClient, error) {
	var kubeManager kubernetes.KubeClient

	if namespace == "" {
		namespace = defaultNamespace
	}

	if c.BaseConfig.OutOfClusterKube {
		var filePath string
		if kubePath == "" {
			glg.Debugf("defaulting to %s", defaultKubeConfig)
			homeDir, err := os.UserHomeDir()
			if err != nil {
				glg.Error(err)
			}

			filePath = filepath.Join(homeDir, defaultKubeConfig)
		} else {
			filePath = kubePath
		}

		cfg, err := ioutil.ReadFile(filePath)
		if err != nil {
			glg.Error("error getting kubeconfig")
		}

		kube, err := kubernetes.NewKubeClient(cfg, namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}

		kubeManager = kube
	} else {
		glg.Debug("creating in cluster kube client")
		kube, err := kubernetes.NewKubeClient(nil, namespace)
		if err != nil {
			glg.Fatal("error creating kubeclient")
		}
		kubeManager = kube
	}

	return &kubeManager, nil
}

func (c *Config) getVaultClient() (*vault.Client, error) {
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

	return &vaultClient, nil
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

func (c *Config) getCert() []byte {
	var cert []byte

	if c.BaseConfig.TestOverwrite {
		glg.Info("trying to read cert file from file")
		rootPath := helpers.OdysseiaRootPath()
		certPath := filepath.Join(rootPath, "eratosthenes", "fixture", "elastic-test-cert.pem")

		cert, _ = ioutil.ReadFile(certPath)

		return cert
	}

	glg.Info("trying to read cert file from pod")
	cert, _ = ioutil.ReadFile(certPathInPod)

	return cert
}
