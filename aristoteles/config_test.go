package aristoteles

import (
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia/aristoteles/configs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAlexandrosConfigCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.AlexandrosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		alexandrosConfig, ok := sut.(*configs.AlexandrosConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, alexandrosConfig.Index)
		assert.NotNil(t, alexandrosConfig.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("ElasticAccessThroughLabel", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvKey, "test")
		cfg := configs.AlexandrosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		alexandrosConfig, ok := sut.(*configs.AlexandrosConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, alexandrosConfig.Index)
		assert.NotNil(t, alexandrosConfig.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvKey)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("OverwriteEnvVariables", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvKey, "minikube")
		os.Setenv(EnvTlSKey, "no")
		cfg := configs.AlexandrosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		alexandrosConfig, ok := sut.(*configs.AlexandrosConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, alexandrosConfig.Index)
		assert.NotNil(t, alexandrosConfig.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvKey)
		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvTlSKey)
	})
}

func TestAnaximanderCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.AnaximanderConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.AnaximanderConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, config.Index)
		assert.Equal(t, 0, config.Created)
		assert.NotNil(t, config.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})
}

func TestDemokritosConfigCreation(t *testing.T) {
	t.Run("StandardConfigCanBeParsed", func(t *testing.T) {
		expected := 0
		searchExpected := "greek"
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.DemokritosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.DemokritosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Elastic)
		assert.Equal(t, expected, config.Created)
		assert.Equal(t, searchExpected, config.SearchWord)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("SearchWordOverwrite", func(t *testing.T) {
		expected := "anotherterm"
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvSearchWord, expected)
		cfg := configs.DemokritosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.DemokritosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Elastic)
		assert.Equal(t, expected, config.SearchWord)
		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvSearchWord)
	})

	t.Run("SearchWordOverwrite", func(t *testing.T) {
		expected := "anotherterm"
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvIndex, expected)
		cfg := configs.DemokritosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.DemokritosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Elastic)
		assert.Equal(t, expected, config.Index)
		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvIndex)
	})
}

func TestDionysiosConfigCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expected := "testrole"
		expectedSecondary := "testsecondaryindex"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvSecondaryIndex, expectedSecondary)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.DionysiosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		dionysiosConfig, ok := sut.(*configs.DionysiosConfig)
		assert.True(t, ok)

		assert.NotNil(t, dionysiosConfig.Elastic)
		assert.NotNil(t, dionysiosConfig.Cache)
		assert.Equal(t, expected, dionysiosConfig.Index)
		assert.Equal(t, expectedSecondary, dionysiosConfig.SecondaryIndex)

		os.Unsetenv(EnvIndex)

		os.Unsetenv(EnvSecondaryIndex)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})
}

func TestDrakonCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expectedIndexes := "testindex;anotherindex"
		expectedRoles := "testrole;anotherrole"
		expectedPodName := "drakonPod"
		os.Setenv(EnvIndexes, expectedIndexes)
		os.Setenv(EnvRoles, expectedRoles)
		os.Setenv(EnvPodName, expectedPodName)
		os.Setenv(EnvHealthCheckOverwrite, "yes")

		cfg := configs.DrakonConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.DrakonConfig)
		assert.True(t, ok)
		assert.Equal(t, defaultNamespace, config.Namespace)
		assert.Equal(t, expectedPodName, config.PodName)
		assert.NotNil(t, config.Elastic)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.Roles)
		assert.NotNil(t, config.Indexes)

		for _, role := range config.Roles {
			assert.Contains(t, expectedRoles, role)
		}

		for _, index := range config.Indexes {
			assert.Contains(t, expectedIndexes, index)
		}

		os.Unsetenv(EnvIndexes)
		os.Unsetenv(EnvRoles)
		os.Unsetenv(EnvPodName)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})
}

func TestHerakleitosCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.HerakleitosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.HerakleitosConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, config.Index)
		assert.Equal(t, 0, config.Created)
		assert.NotNil(t, config.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})
}

func TestHerodotosConfigCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.HerodotosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		alexandrosConfig, ok := sut.(*configs.HerodotosConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, alexandrosConfig.Index)
		assert.NotNil(t, alexandrosConfig.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("ElasticAccessThroughLabel", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvKey, "test")
		cfg := configs.HerodotosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		alexandrosConfig, ok := sut.(*configs.HerodotosConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, alexandrosConfig.Index)
		assert.NotNil(t, alexandrosConfig.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvKey)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("OverwriteEnvVariables", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvKey, "minikube")
		os.Setenv(EnvTlSKey, "no")
		cfg := configs.HerodotosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.HerodotosConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, config.Index)
		assert.NotNil(t, config.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvKey)
		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvTlSKey)
	})
}

func TestParmenidesCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvIndex, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.ParmenidesConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.ParmenidesConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, config.Index)
		assert.Equal(t, 0, config.Created)
		assert.NotNil(t, config.Queue)
		assert.NotNil(t, config.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})
}

func TestPeriandrosCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.PeriandrosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PeriandrosConfig)
		assert.True(t, ok)
		assert.Equal(t, defaultNamespace, config.Namespace)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.HttpClients)

		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("OverwriteForRoleCreation", func(t *testing.T) {
		expectedSolonService := "http://overrwitten:235235"
		expectedPodName := "iamapod-2100-24"
		expectedUsername := "iamapod"
		expectedRoleName := "accessRole"
		expectedIndexNames := "testindex;anotherindex"
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvRole, expectedRoleName)
		os.Setenv(EnvIndex, expectedIndexNames)
		os.Setenv(EnvPodName, expectedPodName)
		os.Setenv(EnvSolonService, expectedSolonService)
		cfg := configs.PeriandrosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PeriandrosConfig)
		assert.True(t, ok)
		assert.Equal(t, defaultNamespace, config.Namespace)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.HttpClients)

		assert.Equal(t, expectedPodName, config.SolonCreationRequest.PodName)
		assert.Equal(t, expectedUsername, config.SolonCreationRequest.Username)
		assert.Equal(t, expectedRoleName, config.SolonCreationRequest.Role)

		for _, index := range config.SolonCreationRequest.Access {
			assert.Contains(t, expectedIndexNames, index)
		}

		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvIndexes)
		os.Unsetenv(EnvRole)
		os.Unsetenv(EnvPodName)
		os.Unsetenv(EnvSolonService)
	})

	t.Run("SecondaryIndex", func(t *testing.T) {
		expectedSolonService := "http://overrwitten:235235"
		expectedPodName := "iamapod-2100-24"
		expectedUsername := "iamapod"
		expectedRoleName := "accessRole"
		expectedIndexName := "testindex"
		expectedSecondayIndex := "anotherindex"
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvRole, expectedRoleName)
		os.Setenv(EnvIndex, expectedIndexName)
		os.Setenv(EnvSecondaryIndex, expectedSecondayIndex)
		os.Setenv(EnvPodName, expectedPodName)
		os.Setenv(EnvSolonService, expectedSolonService)
		cfg := configs.PeriandrosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PeriandrosConfig)
		assert.True(t, ok)
		assert.Equal(t, defaultNamespace, config.Namespace)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.HttpClients)

		assert.Equal(t, expectedPodName, config.SolonCreationRequest.PodName)
		assert.Equal(t, expectedUsername, config.SolonCreationRequest.Username)
		assert.Equal(t, expectedRoleName, config.SolonCreationRequest.Role)

		for _, index := range config.SolonCreationRequest.Access {
			assert.Contains(t, index, "index")
		}

		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvIndexes)
		os.Unsetenv(EnvSecondaryIndex)
		os.Unsetenv(EnvRole)
		os.Unsetenv(EnvPodName)
		os.Unsetenv(EnvSolonService)
	})
}

func TestPeriklesConfigCreation(t *testing.T) {
	t.Run("StandardConfigCanBeParsed", func(t *testing.T) {
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.PeriklesConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PeriklesConfig)
		assert.True(t, ok)
		assert.NotNil(t, config)

		os.Unsetenv(EnvHealthCheckOverwrite)
	})
}

func TestPtolemaiosConfigCreation(t *testing.T) {
	t.Run("StandardConfigCanBeParsed", func(t *testing.T) {
		cfg := configs.PtolemaiosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PtolemaiosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.HttpClients)
		assert.Equal(t, config.Namespace, defaultNamespace)
		assert.False(t, config.RunOnce)
	})

	t.Run("SolonOverwrite", func(t *testing.T) {
		expected := "https://test-solon-service:50232"
		os.Setenv(EnvSolonService, expected)
		cfg := configs.PtolemaiosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PtolemaiosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.HttpClients)
		assert.Equal(t, config.Namespace, defaultNamespace)

		os.Unsetenv(EnvSolonService)
	})

	t.Run("VaultOverwrite", func(t *testing.T) {
		expected := "https://test-vault-service:13241234"
		os.Setenv(EnvVaultService, expected)
		cfg := configs.PtolemaiosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PtolemaiosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Kube)
		assert.Equal(t, config.Namespace, defaultNamespace)
		assert.NotNil(t, config.HttpClients)

		os.Unsetenv(EnvVaultService)
	})

	t.Run("OverWritePodName", func(t *testing.T) {
		expected := "fullpodname-12349235-wkfmf"
		expectedPodName := "fullpodname"
		os.Setenv(EnvPodName, expected)
		cfg := configs.PtolemaiosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PtolemaiosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Kube)
		assert.Equal(t, config.Namespace, defaultNamespace)
		assert.Equal(t, expected, config.FullPodName)
		assert.Equal(t, expectedPodName, config.PodName)

		os.Unsetenv(EnvPodName)
	})

	t.Run("OverRunOnce", func(t *testing.T) {
		os.Setenv(EnvRunOnce, "true")
		cfg := configs.PtolemaiosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PtolemaiosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Kube)
		assert.Equal(t, config.Namespace, defaultNamespace)
		assert.True(t, config.RunOnce)

		os.Unsetenv(EnvRunOnce)
	})
}

func TestSokratesConfigCreation(t *testing.T) {
	t.Run("StandardConfig", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvSearchWord, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.SokratesConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		sokartesConfig, ok := sut.(*configs.SokratesConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, sokartesConfig.SearchWord)
		assert.NotNil(t, sokartesConfig.Elastic)

		os.Unsetenv(EnvSearchWord)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("ElasticAccessThroughLabel", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvSearchWord, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvKey, "test")
		cfg := configs.SokratesConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		sokartesConfig, ok := sut.(*configs.SokratesConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, sokartesConfig.SearchWord)
		assert.NotNil(t, sokartesConfig.Elastic)

		os.Unsetenv(EnvSearchWord)
		os.Unsetenv(EnvKey)
		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("OverwriteEnvVariables", func(t *testing.T) {
		expected := "testrole"
		os.Setenv(EnvSearchWord, expected)
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvKey, "minikube")
		os.Setenv(EnvTlSKey, "no")
		cfg := configs.SokratesConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		sokartesConfig, ok := sut.(*configs.SokratesConfig)
		assert.True(t, ok)
		assert.Equal(t, expected, sokartesConfig.SearchWord)
		assert.NotNil(t, sokartesConfig.Elastic)

		os.Unsetenv(EnvIndex)
		os.Unsetenv(EnvKey)
		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvTlSKey)
	})
}

func TestSolonConfigCreation(t *testing.T) {
	t.Run("StandardConfigCanBeParsed", func(t *testing.T) {
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.SolonConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.SolonConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.Elastic)
		assert.NotNil(t, config.Vault)
		assert.Equal(t, config.Namespace, defaultNamespace)
		assert.Equal(t, config.AccessAnnotation, defaultAccessAnnotation)
		assert.Equal(t, config.RoleAnnotation, defaultRoleAnnotation)

		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("VaultOverwrite", func(t *testing.T) {
		expected := "https://test-vault-service:13241234"
		os.Setenv(EnvVaultService, expected)
		cfg := configs.PtolemaiosConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.PtolemaiosConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Kube)
		assert.NotNil(t, config.Vault)
		assert.Equal(t, config.Namespace, defaultNamespace)

		os.Unsetenv(EnvVaultService)
	})
}

func TestElasticClient(t *testing.T) {
	t.Run("SetElasticService", func(t *testing.T) {
		expectedService := "http://test-service:9200"
		os.Setenv(EnvElasticService, expectedService)
		configManager := Config{
			env:        "test",
			BaseConfig: configs.BaseConfig{},
		}

		sut := configManager.getElasticConfig(false)

		assert.Equal(t, expectedService, sut.Service)
		assert.Equal(t, elasticUsernameDefault, sut.Username)
		assert.Equal(t, elasticPasswordDefault, sut.Password)
		assert.Equal(t, "", sut.ElasticCERT)

		os.Unsetenv(EnvElasticService)
	})

	t.Run("SetElasticUser", func(t *testing.T) {
		expectedUser := "test-user"
		os.Setenv(EnvElasticUser, expectedUser)
		configManager := Config{
			env:        "test",
			BaseConfig: configs.BaseConfig{},
		}

		sut := configManager.getElasticConfig(false)

		assert.Equal(t, elasticServiceDefault, sut.Service)
		assert.Equal(t, expectedUser, sut.Username)
		assert.Equal(t, elasticPasswordDefault, sut.Password)
		assert.Equal(t, "", sut.ElasticCERT)

		os.Unsetenv(EnvElasticUser)
	})

	t.Run("SetElasticPassword", func(t *testing.T) {
		expectedPassword := "test-password"
		os.Setenv(EnvElasticPassword, expectedPassword)
		configManager := Config{
			env:        "test",
			BaseConfig: configs.BaseConfig{},
		}

		sut := configManager.getElasticConfig(false)

		assert.Equal(t, elasticServiceDefault, sut.Service)
		assert.Equal(t, elasticUsernameDefault, sut.Username)
		assert.Equal(t, expectedPassword, sut.Password)
		assert.Equal(t, "", sut.ElasticCERT)

		os.Unsetenv(EnvElasticPassword)
	})

	t.Run("StandardTlsIsRead", func(t *testing.T) {
		configManager := Config{
			env: "test",
			BaseConfig: configs.BaseConfig{
				TestOverwrite: true,
			},
		}

		sut := configManager.getElasticConfig(true)

		assert.Equal(t, elasticServiceDefaultTlS, sut.Service)
		assert.Equal(t, elasticUsernameDefault, sut.Username)
		assert.Equal(t, elasticPasswordDefault, sut.Password)
		assert.NotEqual(t, "", sut.ElasticCERT)
	})

	t.Run("VaultToConf", func(t *testing.T) {
		configManager := Config{
			env: "test",
			BaseConfig: configs.BaseConfig{
				TestOverwrite: true,
			},
		}

		vaultModel := models.ElasticConfigVault{
			Username:    "testuser",
			Password:    "testpassword",
			ElasticCERT: "amegahugecert",
		}

		sut := configManager.mapVaultToConf(&vaultModel, true)

		assert.Equal(t, elasticServiceDefaultTlS, sut.Service)
		assert.Equal(t, vaultModel.Username, sut.Username)
		assert.Equal(t, vaultModel.Password, sut.Password)
		assert.Equal(t, vaultModel.ElasticCERT, sut.ElasticCERT)
	})

	t.Run("TlSDisabledTestConf", func(t *testing.T) {
		configManager := Config{
			env: "somethingelse",
			BaseConfig: configs.BaseConfig{
				TLSEnabled:    false,
				TestOverwrite: true,
				HealthCheck:   false,
			},
		}
		elasticClient, err := configManager.getElasticClient()
		assert.Nil(t, err)
		assert.NotNil(t, elasticClient)
	})

	t.Run("TLSEnabledTestConf", func(t *testing.T) {
		configManager := Config{
			env: "TEST",
			BaseConfig: configs.BaseConfig{
				TLSEnabled:    true,
				TestOverwrite: true,
				HealthCheck:   false,
			},
		}
		elasticClient, err := configManager.getElasticClient()
		assert.Nil(t, err)
		assert.NotNil(t, elasticClient)
	})

}

func TestThalesConfigCreation(t *testing.T) {
	t.Run("StandardConfigCanBeParsed", func(t *testing.T) {
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		cfg := configs.ThalesConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.ThalesConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Elastic)
		assert.Equal(t, defaultChannelName, config.Channel)
		assert.NotNil(t, config.Queue)

		os.Unsetenv(EnvHealthCheckOverwrite)
	})

	t.Run("CanOverwriteChannelName", func(t *testing.T) {
		expected := "testChannel"
		os.Setenv(EnvHealthCheckOverwrite, "yes")
		os.Setenv(EnvChannel, expected)
		cfg := configs.ThalesConfig{}

		sut, err := NewConfig(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, sut)

		config, ok := sut.(*configs.ThalesConfig)
		assert.True(t, ok)
		assert.NotNil(t, config.Elastic)
		assert.Equal(t, expected, config.Channel)

		os.Unsetenv(EnvHealthCheckOverwrite)
		os.Unsetenv(EnvChannel)
	})
}
