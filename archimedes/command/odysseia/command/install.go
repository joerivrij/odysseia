package command

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/command"
	elastic "github.com/odysseia/archimedes/command/kubernetes/command"
	"github.com/odysseia/archimedes/util"
	"github.com/odysseia/plato/generator"
	"github.com/odysseia/plato/harbor"
	"github.com/odysseia/plato/helm"
	"github.com/odysseia/plato/kubernetes"
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

func Install() *cobra.Command {
	var (
		namespace        string
		kubePath         string
		themistoklesPath string
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

			if themistoklesPath == "" {
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
				themistoklesPath = filepath.Join(l, "themistokles", "odysseia", "charts")
			}

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

			glg.Info("creating a new install for odysseia")

			odysseia := OdysseiaInstaller{
				Namespace:        namespace,
				ConfigPath:       "",
				CurrentPath:      "",
				ThemistoklesRoot: themistoklesPath,
				Charts:           Themistokles{},
				Config:           CurrentInstallConfig{},
				Kube:             kubeManager,
				Helm:             helmManager,
				Harbor:           nil,
			}

			odysseia.installOdysseiaComplete()
		},
	}
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&kubePath, "kubepath", "k", "", "kubeconfig filepath defaults to ~/.kube/config")
	cmd.PersistentFlags().StringVarP(&themistoklesPath, "themistokles", "t", "", "the path to your helm chart")

	return cmd
}

type OdysseiaInstaller struct {
	Namespace        string
	ConfigPath       string
	CurrentPath      string
	ThemistoklesRoot string
	Charts           Themistokles
	Kube             kubernetes.KubeClient
	Helm             helm.HelmClient
	Harbor           harbor.Client
	Config           CurrentInstallConfig
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

type CurrentInstallConfig struct {
	ElasticPassword string `yaml:"elastic-password"`
	HarborPassword  string `yaml:"harbor-password"`
}

func (o *OdysseiaInstaller) installOdysseiaComplete() {
	err := o.preSteps()
	if err != nil {
		glg.Error("error during setup phase")
		glg.Fatal(err)
	}

	err = o.fillHelmChartPaths()
	if err != nil {
		glg.Error("error during setup phase")
		glg.Fatal(err)
	}

	//1. install elastic
	err = o.createElastic()
	if err != nil {
		glg.Fatal(err)
	}

	//2. install harbor
	err = o.installHarborHelmChart()
	if err != nil {
		glg.Fatal(err)
	}

	//2.a create harbor project etc.
	err = o.setupHarbor()
	if err != nil {
		glg.Fatal(err)
	}

	glg.Infof("created harbor project %s at %s", command.DefaultNamespace, command.DefaultHarborUrl)

	//3. create&push images
	//5. install vault
	//6. install solon
	//7. install app

	//save config to configpath
	currentConfig, err := yaml.Marshal(o.Config)
	if err != nil {
		glg.Error(err)
	}
	currentConfigPath := filepath.Join(o.ConfigPath, "config.yaml")
	util.WriteFile(currentConfig, currentConfigPath)
	// copy everything to currentdir
	o.copyToCurrentDir()
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

	rls, err := o.Helm.Install(o.Charts.ElasticSearch)
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
	crt, key, _ := generator.GenerateKeyAndCertSet(hosts, org)
	certData := make(map[string][]byte)
	certData["tls.key"] = key
	certData["tls.crt"] = crt

	secretName := command.DefaultHarborCertSecretName

	err = o.Kube.Configuration().CreateSecret(command.DefaultNamespace, secretName, certData)
	if err != nil {
		return err
	}

	values := map[string]interface{}{
		"harborAdminPassword":          password,
		"expose.tls.secret.secretName": secretName,
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

func (o *OdysseiaInstaller) setupHarbor() error {
	//wait for harbor to install
	ticker := time.NewTicker(time.Second)
	timeout := time.After(60 * time.Second)
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
				if strings.Contains(pod.Name, "harbor") {
					if pod.Status.Phase != "Running" {
						glg.Infof("pod: %s not ready", pod.Name)
						podsReady = false
					}

					readyStatusFound := false
					for _, condition := range pod.Status.Conditions {
						if condition.Type == "Ready" {
							readyStatusFound = true
							break
						}
					}

					if !readyStatusFound {
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
		return fmt.Errorf("harbor pods have not become healthy after 60 seconds")
	}

	err := o.Harbor.CreateProject(command.DefaultNamespace, false)
	return err
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
		glg.Info("directory already exists")
	}

	currentDir := filepath.Join(configDir, "current")

	o.CurrentPath = currentDir

	if _, err := os.Stat(currentDir); os.IsNotExist(err) {
		// create dir
		glg.Infof("creating current config at: %s", currentDir)
		os.Mkdir(currentDir, 0755)
	} else {
		glg.Debug("current dir already exists")
	}

	return nil
}
