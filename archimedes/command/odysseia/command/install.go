package command

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/command"
	imageCommand "github.com/odysseia/archimedes/command/images/command"
	elastic "github.com/odysseia/archimedes/command/kubernetes/command"
	vaultCommand "github.com/odysseia/archimedes/command/vault/command"
	"github.com/odysseia/archimedes/util"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/certificates"
	"github.com/odysseia/plato/generator"
	"github.com/odysseia/plato/harbor"
	"github.com/odysseia/plato/helm"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	configFilePath string = ".odysseia"
)

//go:embed "config"
var configPath embed.FS

//go:embed "apps"
var whiteListFile embed.FS

func Install() *cobra.Command {
	var (
		namespace        string
		kubePath         string
		themistoklesPath string
		odysseiaRootPath string
		profile          string
		tag              string
		destinationRepo  string
		unseal           bool
		build            bool
	)
	cmd := &cobra.Command{
		Use:   "install",
		Short: "parse a list of words",
		Long: `Allows you to parse a list of words to be used by demokritos
- Filepath
`,
		Run: func(cmd *cobra.Command, args []string) {

			if namespace == "" {
				glg.Debugf("defaulting to %s", command.DefaultNamespace)
				namespace = command.DefaultNamespace
			}

			if kubePath == "" {
				glg.Debugf("defaulting to %s", command.DefaultKubeConfig)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					glg.Error(err)
				}

				kubePath = filepath.Join(homeDir, command.DefaultKubeConfig)
			}

			if odysseiaRootPath == "" {
				_, callingFile, _, _ := runtime.Caller(0)
				callingDir := filepath.Dir(callingFile)
				dirParts := strings.Split(callingDir, string(os.PathSeparator))
				var odysseiaPath []string
				for i, part := range dirParts {
					if part == "odysseia" {
						odysseiaPath = dirParts[0 : i+1]
						break
					}
				}
				l := "/"
				for _, path := range odysseiaPath {
					l = filepath.Join(l, path)
				}

				odysseiaRootPath = l
			}

			themistoklesPath = filepath.Join(odysseiaRootPath, "themistokles", "odysseia", "charts")

			cfg, err := ioutil.ReadFile(kubePath)
			if err != nil {
				glg.Error("error getting kubeconfig")
			}

			kubeManager, err := kubernetes.NewKubeClient(cfg, namespace)
			if err != nil {
				glg.Fatal("error creating kubeclient")
			}

			helmManager, err := helm.NewHelmClient(cfg, namespace)
			if err != nil {
				glg.Fatal("error creating helmclient")
			}

			glg.Info("getting config from yaml files")

			config, err := configPath.ReadFile(fmt.Sprintf("config/%s.yaml", profile))
			if err != nil {
				glg.Fatal("error reading config files")
			}

			var valueOverwrite configs.ValueOverwrite
			err = yaml.Unmarshal(config, &valueOverwrite)
			if err != nil {
				glg.Fatal("error marshalling yaml")
			}

			if tag == "" {
				gitTag, err := util.ExecCommandWithReturn(`git rev-parse --short HEAD`, odysseiaRootPath)
				if err != nil {
					glg.Fatal("error getting gitref")
				}

				tag = gitTag
			}

			glg.Info("creating a new install for odysseia")

			whiteList, err := whiteListFile.ReadFile("apps/whitelist.yaml")
			if err != nil {
				glg.Fatal("error reading config files")
			}

			var wl models.WhiteList
			err = yaml.Unmarshal(whiteList, &wl)
			if err != nil {
				glg.Fatal("error marshalling yaml")
			}

			newConfig := models.CurrentInstallConfig{
				ElasticPassword: "",
				HarborPassword:  "",
				VaultRootToken:  "",
				VaultUnsealKey:  "",
			}

			odysseia := OdysseiaInstaller{
				Namespace:           namespace,
				ConfigPath:          "",
				CurrentPath:         "",
				ThemistoklesRoot:    themistoklesPath,
				OdysseiaRoot:        odysseiaRootPath,
				Repo:                destinationRepo,
				Charts:              Themistokles{},
				Config:              newConfig,
				ValueConfig:         valueOverwrite,
				Kube:                kubeManager,
				Helm:                helmManager,
				GitTag:              tag,
				Profile:             profile,
				Harbor:              nil,
				ChartsToInstall:     []string{},
				WhiteList:           wl.AppsToInstall,
				PrivateImagesInRepo: false,
				ExternalRepo:        true,
				Build:               build,
			}

			err = odysseia.installOdysseiaComplete(unseal)
			if err != nil {
				glg.Error(err)
				os.Exit(1)
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&kubePath, "kubepath", "k", "", "kubeconfig filepath defaults to ~/.kube/config")
	cmd.PersistentFlags().StringVarP(&odysseiaRootPath, "odysseia", "o", "", "the path to your odysseia root")
	cmd.PersistentFlags().BoolVarP(&unseal, "unseal", "u", false, "whether to unseal vault")
	cmd.PersistentFlags().StringVarP(&profile, "profile", "p", "docker-desktop", "what profile to use for the install")
	cmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "image tag")
	cmd.PersistentFlags().StringVarP(&destinationRepo, "dest", "d", "core.harbor.domain:30003/odysseia", "destination repo address")
	cmd.PersistentFlags().BoolVarP(&build, "build", "b", true, "whether to build images")

	return cmd
}

type OdysseiaInstaller struct {
	Namespace           string
	ConfigPath          string
	CurrentPath         string
	ThemistoklesRoot    string
	OdysseiaRoot        string
	GitTag              string
	Profile             string
	Repo                string
	PrivateImagesInRepo bool
	ExternalRepo        bool
	Build               bool
	ChartsToInstall     []string
	WhiteList           []string
	Charts              Themistokles
	Kube                kubernetes.KubeClient
	Helm                helm.HelmClient
	Harbor              harbor.Client
	Config              models.CurrentInstallConfig
	ValueConfig         configs.ValueOverwrite
}

type Themistokles struct {
	ElasticSearch string
	Vault         string
	Solon         string
	Harbor        string
	Kibana        string
	Docs          string
	Tests         string
	Apis          []string
}

func (o *OdysseiaInstaller) installOdysseiaComplete(unseal bool) error {
	err := o.preSteps()
	if err != nil {
		return err
	}

	err = o.fillHelmChartPaths()
	if err != nil {
		return err
	}

	err = o.filterHelmChartsToBeInstalled()
	if err != nil {
		return err
	}

	defer func() {
		//save config to configpath
		err = o.checkConfigForEmpty()
		if err != nil {
			glg.Error(err)
		}

		currentConfig, err := yaml.Marshal(o.Config)
		if err != nil {
			glg.Error(err)
		}

		glg.Info(string(currentConfig))
		currentConfigPath := filepath.Join(o.ConfigPath, "config.yaml")
		util.WriteFile(currentConfig, currentConfigPath)
		// copy everything to currentdir
		err = o.copyToCurrentDir()
		if err != nil {
			glg.Error(err)
		}
	}()

	//1. install elastic
	installElastic := false
	for _, install := range o.ChartsToInstall {
		if install == "elasticsearch" {
			installElastic = true
			break
		}
	}

	if installElastic {
		err = o.createElastic()
		if err != nil {
			return err
		}
	}

	installHarbor := false
	for _, install := range o.ChartsToInstall {
		if install == "harbor" {
			installHarbor = true
			break
		}
	}

	if installHarbor {
		//2. install harbor
		err = o.installHarborHelmChart()
		if err != nil {
			return err
		}

		//2.a create harbor project etc.
		err = o.setupHarbor()
		if err != nil {
			return err
		}

		glg.Infof("created harbor project %s at %s", command.DefaultNamespace, command.DefaultHarborUrl)

		//2.b. docker login
		err = o.dockerLogin()
		if err != nil {
			return err
		}

		o.PrivateImagesInRepo = true
	}

	installVault := false
	for _, install := range o.ChartsToInstall {
		if install == "vault" {
			installVault = true
			break
		}
	}

	if installVault {
		//4. install vault
		err = o.installVaultHelmChart()
		if err != nil {
			return err
		}

		//4b. provision vault
		err = o.setupVault()
		if err != nil {
			return err
		}
	}

	if unseal {
		vaultCommand.UnsealVault("", o.Namespace, o.Kube)
	}

	var repoName string

	if string(o.Repo[len(o.Repo)-1]) != "/" {
		repoName = o.Repo + "/"
	} else {
		repoName = o.Repo
	}

	image := map[string]interface{}{
		"tag":       o.GitTag,
		"imageRepo": repoName,
	}

	if o.PrivateImagesInRepo {
		var pullSecret string
		if o.ValueConfig.Harbor.Expose.TLS.Secret.SecretName == "" {
			pullSecret = command.DefaultHarborCertSecretName
		} else {
			pullSecret = o.ValueConfig.Harbor.Expose.TLS.Secret.SecretName
		}
		image["pullSecret"] = pullSecret
	}

	config := map[string]interface{}{
		"externalRepo":        o.ExternalRepo,
		"privateImagesInRepo": o.PrivateImagesInRepo,
	}

	if o.ExternalRepo {
		config["pullPolicy"] = "Always"
	}

	values := map[string]interface{}{
		"images": image,
		"config": config,
	}

	installSolon := false
	for _, install := range o.ChartsToInstall {
		if install == "solon" {
			installSolon = true
			break
		}
	}

	if installSolon {
		//6. install solon
		splitName := strings.Split(o.Charts.Solon, "/")
		chartName := strings.ToLower(splitName[len(splitName)-1])
		err = o.installHelmChart(chartName, o.Charts.Solon, values)
		if err != nil {
			return err
		}
	}

	//7. install app
	err = o.installAppsHelmChart(values)
	if err != nil {
		return err
	}

	installTests := false
	for _, install := range o.ChartsToInstall {
		if install == "tests" {
			installTests = true
			break
		}
	}

	//8. install tests
	if installTests {
		splitName := strings.Split(o.Charts.Tests, "/")
		chartName := strings.ToLower(splitName[len(splitName)-1])
		err = o.installHelmChart(chartName, o.Charts.Tests, values)
		if err != nil {
			return err
		}
	}

	installDocs := false
	for _, install := range o.ChartsToInstall {
		if install == "docs" {
			installDocs = true
			break
		}
	}

	//9. install docs
	if installDocs {
		splitName := strings.Split(o.Charts.Docs, "/")
		chartName := strings.ToLower(splitName[len(splitName)-1])
		err = o.installHelmChart(chartName, o.Charts.Docs, values)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *OdysseiaInstaller) copyToCurrentDir() error {
	files, err := ioutil.ReadDir(o.ConfigPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcDir := filepath.Join(o.ConfigPath, f.Name())
		destDir := filepath.Join(o.CurrentPath, f.Name())

		err := util.CopyFileContents(srcDir, destDir)
		if err != nil {
			glg.Error(err)
		}
	}

	return nil
}

func (o *OdysseiaInstaller) createElastic() error {
	p12Path, err := elastic.CreateElasticP12(o.Kube, o.Namespace, o.ConfigPath)
	if err != nil {
		return err
	}

	p12File, err := os.ReadFile(p12Path)
	if err != nil {
		return err
	}

	pemDst := filepath.Join(o.ConfigPath, "elastic-certificate.pem")
	cmd := fmt.Sprintf(`openssl pkcs12 -nodes -passin pass:'' -in %s -out %s`, p12Path, pemDst)

	err = util.ExecCommand(cmd, "/")
	if err != nil {
		return err
	}

	pemFile, err := os.ReadFile(pemDst)
	if err != nil {
		return err
	}
	crtDst := filepath.Join(o.ConfigPath, "elastic-certificate.crt")
	crtFile, err := elastic.GenerateCrtFromPem(pemFile)
	if err != nil {
		return err
	}

	util.WriteFile(crtFile, crtDst)

	//create secrets
	glg.Info("certs for ES tls mode generated applying them as secrets")

	secretNameP12 := "elastic-certificates"
	dataP12 := make(map[string][]byte)
	dataP12["elastic-certificates.p12"] = p12File

	err = o.Kube.Configuration().CreateSecret(o.Namespace, secretNameP12, dataP12)
	if err != nil {
		return err
	}

	secretNamePem := "elastic-certificate-pem"
	dataPem := make(map[string][]byte)
	dataPem["elastic-certificate.pem"] = pemFile

	err = o.Kube.Configuration().CreateSecret(o.Namespace, secretNamePem, dataPem)
	if err != nil {
		return err
	}

	secretNameCrt := "elastic-certificate-crt"
	dataCrt := make(map[string][]byte)
	dataCrt["elastic-certificate.crt"] = crtFile

	err = o.Kube.Configuration().CreateSecret(o.Namespace, secretNameCrt, dataCrt)
	if err != nil {
		return err
	}

	//create elastic login
	password, err := generator.RandomPassword(24)
	if err != nil {
		return err
	}

	o.Config.ElasticPassword = password

	glg.Debug(password)

	data := make(map[string][]byte)
	data["password"] = []byte(password)
	data["username"] = []byte("elastic")

	err = o.Kube.Configuration().CreateSecret(o.Namespace, command.DefaultSecretName, data)
	if err != nil {
		return err
	}

	glg.Infof("created secret with name %s", command.DefaultSecretName)

	values, err := o.parseValueOverwrite("elastic")
	if err != nil {
		return err
	}

	rls, err := o.Helm.InstallWithValues(o.Charts.ElasticSearch, values)
	if err != nil {
		return err
	}
	glg.Info(rls.Name)

	return nil
}

func (o *OdysseiaInstaller) installHarborHelmChart() error {
	//create harbor login
	password, err := generator.RandomPassword(24)
	if err != nil {
		return err
	}

	o.Config.HarborPassword = password

	glg.Debug(password)

	data := make(map[string]string)
	data["docker-server"] = "harbor"
	data["docker-username"] = command.DefaultAdmin
	data["docker-password"] = password
	data["docker-email"] = "odysseia@example.com"

	secret, _ := json.Marshal(data)
	secretData := map[string]string{".dockerconfigjson": string(secret)}

	err = o.Kube.Configuration().CreateDockerSecret(o.Namespace, command.DefaultDockerRegistrySecret, secretData)
	if err != nil {
		return err
	}

	hosts := []string{
		command.DefaultHarbor,
	}
	org := []string{
		command.DefaultNamespace,
	}

	validity := 3650

	certClient, err := certificates.NewCertGeneratorClient(org, validity)
	err = certClient.InitCa()
	if err != nil {
		return err
	}

	crt, key, _ := certClient.GenerateKeyAndCertSet(hosts, validity)
	certData := make(map[string][]byte)
	certData["tls.key"] = key
	certData["tls.crt"] = crt

	secretName := command.DefaultHarborCertSecretName

	err = o.Kube.Configuration().CreateSecret(command.DefaultNamespace, secretName, certData)
	if err != nil {
		return err
	}

	o.ValueConfig.Harbor.HarborAdminPassword = password
	o.ValueConfig.Harbor.Expose.TLS.Secret.SecretName = secretName

	values, err := o.parseValueOverwrite("harbor")
	if err != nil {
		return err
	}

	rls, err := o.Helm.InstallWithValues(o.Charts.Harbor, values)
	if err != nil {
		return err
	}
	glg.Info(rls.Name)

	harborManager, _ := harbor.NewHarborClient(command.DefaultHarborUrl, command.DefaultAdmin, password, crt)
	o.Harbor = harborManager

	return nil
}

func (o *OdysseiaInstaller) installVaultHelmChart() error {
	rls, err := o.Helm.Install(o.Charts.Vault)
	if err != nil {
		return err
	}
	glg.Info(rls.Name)

	return nil
}

func (o *OdysseiaInstaller) installAppsHelmChart(valuesOverwrite map[string]interface{}) error {
	//wait for solon to install
	ticker := time.NewTicker(5 * time.Second)
	timeout := time.After(360 * time.Second)
	var ready bool
	for {
		select {
		case <-ticker.C:
			pods, err := o.Kube.Workload().List(command.DefaultNamespace)
			if err != nil {
				continue
			}

			podsReady := true
			for _, pod := range pods.Items {
				if !podsReady {
					break
				}
				if strings.Contains(pod.Name, "solon") {
					if pod.Status.Phase != "Running" {
						glg.Infof("pod: %s not ready", pod.Name)
						podsReady = false
					}

					readyStatusFound := false
					for _, condition := range pod.Status.Conditions {
						if condition.Type == "Ready" && condition.Status == "True" {
							readyStatusFound = true
							break
						}
					}

					if !readyStatusFound {
						glg.Infof("pod: %s not ready", pod.Name)
						podsReady = false
					}
				}
			}

			if podsReady {
				ready = true
				ticker.Stop()
			} else {
				continue
			}

		case <-timeout:
			glg.Error("timed out")
			ticker.Stop()
		}
		break
	}

	if !ready {
		return fmt.Errorf("solon pod has not become healthy after 120 seconds")
	}

	for _, chart := range o.Charts.Apis {
		for _, install := range o.ChartsToInstall {
			if strings.Contains(chart, install) {
				splitName := strings.Split(chart, "/")
				chartName := strings.ToLower(splitName[len(splitName)-1])
				err := o.installHelmChart(chartName, chart, valuesOverwrite)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (o *OdysseiaInstaller) installHelmChart(name, chartPath string, valuesOverwrite map[string]interface{}) error {
	if o.Build {
		imageCommand.BuildImageSet(o.OdysseiaRoot, name, o.GitTag, o.Repo, true)
	}

	rls, err := o.Helm.InstallWithValues(chartPath, valuesOverwrite)
	if err != nil {
		return err
	}
	glg.Info(rls.Name)

	return nil
}

func (o *OdysseiaInstaller) setupVault() error {
	//wait for vault to install
	ticker := time.NewTicker(time.Second)
	timeout := time.After(120 * time.Second)
	var ready bool
	for {
		select {
		case <-ticker.C:
			pods, err := o.Kube.Workload().List(command.DefaultNamespace)
			if err != nil {
				continue
			}

			podsReady := true
			for _, pod := range pods.Items {
				if !podsReady {
					break
				}
				if strings.Contains(pod.Name, "vault") {
					if pod.Status.Phase != "Running" {
						glg.Infof("pod: %s not ready", pod.Name)
						podsReady = false
					}
				}
			}

			if podsReady {
				ready = true
				ticker.Stop()
			} else {
				continue
			}

		case <-timeout:
			glg.Error("timed out")
			ticker.Stop()
		}
		break
	}

	if !ready {
		return fmt.Errorf("vault pod has not started after 30 seconds")
	}

	time.Sleep(1000 * time.Millisecond)
	vaultConfig, err := vaultCommand.NewVaultFlow(o.Namespace, o.Kube)
	o.Config.VaultUnsealKey = vaultConfig.UnsealKeysHex[0]
	o.Config.VaultRootToken = vaultConfig.RootToken

	return err
}

func (o *OdysseiaInstaller) setupHarbor() error {
	//wait for harbor to install
	ticker := time.NewTicker(time.Second)
	timeout := time.After(120 * time.Second)
	var ready bool
	for {
		select {
		case <-ticker.C:
			pods, err := o.Kube.Workload().List(command.DefaultNamespace)
			if err != nil {
				continue
			}

			podsReady := true
			for _, pod := range pods.Items {
				if !podsReady {
					break
				}
				if strings.Contains(pod.Name, "harbor-core") {
					if pod.Status.Phase != "Running" {
						glg.Infof("pod: %s not ready", pod.Name)
						podsReady = false
					}

					readyStatusFound := false
					for _, condition := range pod.Status.Conditions {
						if condition.Type == "Ready" && condition.Status == "True" {
							readyStatusFound = true
							break
						}
					}

					if !readyStatusFound {
						glg.Infof("pod: %s not ready", pod.Name)
						podsReady = false
					}
				}
			}

			if podsReady {
				ready = true
				ticker.Stop()
			} else {
				continue
			}

		case <-timeout:
			glg.Error("timed out")
			ticker.Stop()
		}
		break
	}

	if !ready {
		return fmt.Errorf("harbor pods have not become healthy after 120 seconds")
	}

	err := o.Harbor.CreateProject(command.DefaultNamespace, true)
	return err
}

func (o *OdysseiaInstaller) dockerLogin() error {
	dockerCommand := fmt.Sprintf("docker login %s --username %s --password %s", o.ValueConfig.Harbor.ExternalURL, command.DefaultAdmin, o.ValueConfig.Harbor.HarborAdminPassword)
	err := util.ExecCommand(dockerCommand, "/")
	if err != nil {
		return err
	}
	return nil
}

func (o *OdysseiaInstaller) filterHelmChartsToBeInstalled() error {
	releases, err := o.Helm.List()
	if err != nil {
		return err
	}

	var toInstall []string

	for _, whitelist := range o.WhiteList {
		found := false
		for _, release := range releases {
			if whitelist == release.Name {
				found = true
			}
		}

		if !found {
			toInstall = append(toInstall, whitelist)
		}
	}

	o.ChartsToInstall = toInstall

	return nil
}

func (o *OdysseiaInstaller) fillHelmChartPaths() error {
	files, err := ioutil.ReadDir(o.ThemistoklesRoot)
	if err != nil {
		return err
	}

	for _, f := range files {
		elements := reflect.ValueOf(&o.Charts).Elem()
		found := false
		for i := 0; i < elements.NumField(); i++ {
			fieldName := elements.Type().Field(i).Name
			if strings.ToLower(fieldName) == f.Name() {
				path := filepath.Join(o.ThemistoklesRoot, f.Name())
				elements.FieldByName(fieldName).SetString(path)
				found = true
			}
		}

		if !found {
			path := filepath.Join(o.ThemistoklesRoot, f.Name())
			o.Charts.Apis = append(o.Charts.Apis, path)
		}
	}

	return nil
}

func (o *OdysseiaInstaller) preSteps() error {
	err := o.Kube.Namespaces().Create(o.Namespace)
	if err != nil {
		return err
	}

	//create config file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, configFilePath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// create dir
		glg.Infof("config directory does not exist... creating at: %s", configDir)
		os.Mkdir(configDir, 0755)
	}

	t := time.Now()
	formattedDate := fmt.Sprintf("%d%02d%02d",
		t.Year(), t.Month(), t.Day())
	newFullInstallDir := filepath.Join(configDir, formattedDate)

	o.ConfigPath = newFullInstallDir

	if _, err := os.Stat(newFullInstallDir); os.IsNotExist(err) {
		// create dir
		glg.Infof("creating new config at: %s", newFullInstallDir)
		os.Mkdir(newFullInstallDir, 0755)
	} else {
		glg.Infof("directory: %s already exists", newFullInstallDir)
	}

	currentDir := filepath.Join(configDir, "current")

	o.CurrentPath = currentDir

	if _, err := os.Stat(currentDir); os.IsNotExist(err) {
		// create dir
		glg.Infof("creating current config at: %s", currentDir)
		os.Mkdir(currentDir, 0755)
	} else {
		glg.Infof("directory: %s already exists", currentDir)
	}

	return nil
}

func (o *OdysseiaInstaller) parseValueOverwrite(service string) (map[string]interface{}, error) {
	var unmarshalledFields map[string]interface{}

	switch service {
	case "harbor":
		harborValues, err := yaml.Marshal(o.ValueConfig.Harbor)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(harborValues, &unmarshalledFields)
	case "elastic":
		elasticValues, err := yaml.Marshal(o.ValueConfig.Elastic)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(elasticValues, &unmarshalledFields)
	}

	return unmarshalledFields, nil
}

func (o *OdysseiaInstaller) checkConfigForEmpty() error {
	currentConfigPath := filepath.Join(o.CurrentPath, "config.yaml")
	fromFile, err := os.ReadFile(currentConfigPath)
	if err != nil {
		return err
	}

	var currentConfig models.CurrentInstallConfig
	err = yaml.Unmarshal(fromFile, &currentConfig)
	if err != nil {
		return err
	}

	if o.Config.HarborPassword == "" {
		o.Config.HarborPassword = currentConfig.HarborPassword
	}

	if o.Config.ElasticPassword == "" {
		o.Config.ElasticPassword = currentConfig.ElasticPassword
	}

	if o.Config.VaultRootToken == "" {
		o.Config.VaultRootToken = currentConfig.VaultRootToken
	}

	if o.Config.VaultUnsealKey == "" {
		o.Config.VaultUnsealKey = currentConfig.VaultUnsealKey
	}

	return nil
}
