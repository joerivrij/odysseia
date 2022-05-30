package aristoteles

import (
	"embed"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/cache"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/queue"
	"github.com/odysseia/plato/service"
	"github.com/odysseia/plato/vault"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	env        string
	BaseConfig configs.BaseConfig
}

//go:embed base
var base embed.FS

func NewConfig(v interface{}) (interface{}, error) {
	glg.Info("ARISTOTELES")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"Τριών δει παιδεία: φύσεως, μαθήσεως, ασκήσεως.\"")
	glg.Info("\"Education needs these three: natural endowment, study, practice.\"")
	glg.Info(strings.Repeat("~", 37))

	env := os.Getenv(EnvKey)
	if env == "" {
		env = "LOCAL"
	}

	envDir := strings.ToLower(env)

	glg.Infof("getting config: %s from yaml files", envDir)

	config, err := base.ReadFile(fmt.Sprintf("%s/%s/%s", baseDir, envDir, configFileName))
	if err != nil {
		return nil, err
	}

	var baseConfig configs.BaseConfig
	err = yaml.Unmarshal(config, &baseConfig)
	if err != nil {
		return nil, err
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

	client := service.NewHttpClient()
	baseConfig.HttpClient = client

	newConfig := Config{BaseConfig: baseConfig}
	newConfig.env = env

	glg.Infof("env set to: %s", newConfig.env)

	healthCheck := true
	healthCheckOverwrite := os.Getenv(EnvHealthCheckOverwrite)
	if healthCheckOverwrite == "yes" || healthCheckOverwrite == "true" {
		healthCheck = false
	}

	sidecarOverwrite := false
	envSidecarOverwrite := os.Getenv(EnvSidecarOverwrite)
	if envSidecarOverwrite == "yes" || envSidecarOverwrite == "true" {
		sidecarOverwrite = true
	}

	newConfig.BaseConfig.SidecarOverwrite = sidecarOverwrite

	if newConfig.BaseConfig.HealthCheckOverwrite {
		healthCheck = !newConfig.BaseConfig.HealthCheckOverwrite
	}

	newConfig.BaseConfig.HealthCheck = healthCheck

	glg.Info("validating config")
	valid, err := newConfig.configValidator(v)
	if !valid || err != nil {
		return nil, err
	}

	glg.Info("config is valid")
	glg.Info("setting remaining fields")
	readConfig, err := newConfig.New(v)
	if err != nil {
		return nil, err
	}

	glg.Info("config created returning")

	return readConfig, nil
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
			case "SecondaryIndex":
				indexName := c.getStringFromEnv(EnvSecondaryIndex, c.BaseConfig.Index)
				e.FieldByName(fieldName).SetString(indexName)
			case "Job":
				indexName := c.getStringFromEnv(EnvJobName, defaultJobName)
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
			case "Channel":
				vs := c.getStringFromEnv(EnvChannel, defaultChannelName)
				e.FieldByName(fieldName).SetString(vs)
			case "SearchWord":
				vs := c.getStringFromEnv(EnvSearchWord, defaultSearchWord)
				e.FieldByName(fieldName).SetString(vs)
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
		case reflect.TypeOf((*elastic.Client)(nil)).Elem():
			k, err := c.getElasticClient()
			if err != nil {
				glg.Fatal("error getting kubeconfig")
			}
			kv := reflect.ValueOf(k)
			e.FieldByName(fieldName).Set(kv)

		case reflect.TypeOf((*kubernetes.KubeClient)(nil)).Elem():
			k, err := c.getKubeClient()
			if err != nil {
				glg.Fatal("error getting kubeconfig")
			}
			kv := reflect.ValueOf(k)
			e.FieldByName(fieldName).Set(kv)

		case reflect.TypeOf((*vault.Client)(nil)).Elem():
			reflectedVault, err := c.getVaultClient()
			if err != nil {
				glg.Fatal("error getting vaultClient")
			}
			vv := reflect.ValueOf(reflectedVault)
			e.FieldByName(fieldName).Set(vv)

		case reflect.TypeOf((*queue.Client)(nil)).Elem():
			mq, err := c.getMqQueueClient()
			if err != nil {
				glg.Fatal("error getting kubemq client")
			}
			mv := reflect.ValueOf(mq)
			e.FieldByName(fieldName).Set(mv)

		case reflect.TypeOf((*cache.Client)(nil)).Elem():
			reflectedBadger, err := c.getBadgerClient()
			if err != nil {
				glg.Fatal("error getting badgerClient")
			}
			vv := reflect.ValueOf(reflectedBadger)
			e.FieldByName(fieldName).Set(vv)

		case reflect.TypeOf((*service.OdysseiaClient)(nil)).Elem():
			mq, err := c.getOdysseiaClient()
			if err != nil {
				glg.Fatal("error getting odysseia client")
			}
			mv := reflect.ValueOf(mq)
			e.FieldByName(fieldName).Set(mv)

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
	secondaryAccess := c.getStringFromEnv(EnvSecondaryIndex, "")
	if secondaryAccess != "" {
		envAccess = append(envAccess, secondaryAccess)
	}
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
