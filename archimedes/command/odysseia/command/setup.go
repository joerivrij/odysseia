package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/command"
	"github.com/odysseia/archimedes/util"
	"github.com/odysseia/plato/helm"
	"github.com/odysseia/plato/kubernetes"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	NginxRepoPath  string = "https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.0.19/ingress-nginx-4.0.19.tgz"
	NginxNamespace string = "ingress-nginx"
	QueueName      string = "queue"
)

func Setup() *cobra.Command {
	var (
		namespace string
		kubePath  string
		profile   string
	)
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup an empty cluster",
		Long: `When starting a new cluster create a new 
- Filepath
`,
		Run: func(cmd *cobra.Command, args []string) {
			glg.Green("setting up")

			var odysseiaRootPath string
			var queuePath string

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

			themistoklesPath := filepath.Join(odysseiaRootPath, "themistokles", "odysseia", "charts")

			files, err := ioutil.ReadDir(themistoklesPath)
			if err != nil {
				glg.Error(err)
			}

			for _, f := range files {
				if f.Name() == QueueName {
					queuePath = filepath.Join(themistoklesPath, f.Name())
					break
				}
			}

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

			odysseia := OdysseiaSetup{
				Namespace: namespace,
				Kube:      kubeManager,
				Helm:      helmManager,
				Profile:   profile,
				QueuePath: queuePath,
			}

			err = odysseia.firstTimeSetup()
			if err != nil {
				glg.Error(err)
				os.Exit(1)
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&kubePath, "kubepath", "k", "", "kubeconfig filepath defaults to ~/.kube/config")
	cmd.PersistentFlags().StringVarP(&profile, "profile", "p", "docker-desktop", "what profile to use for the install")

	return cmd
}

type OdysseiaSetup struct {
	QueuePath string
	Profile   string
	Namespace string
	Kube      kubernetes.KubeClient
	Helm      helm.HelmClient
}

func (o *OdysseiaSetup) firstTimeSetup() error {
	// creates if the namespace does not exist does nothing if it is found
	err := o.Kube.Namespaces().Create(o.Namespace)
	if err != nil {
		return err
	}

	switch o.Profile {
	case "docker-desktop":
		rls, err := o.Helm.InstallNamespaced(NginxRepoPath, NginxNamespace, true)
		if err != nil {
			return err
		}

		glg.Debugf("created nginx release on docker-desktop %v in ns %v", rls.Name, rls.Namespace)
	case "do":
		rls, err := o.Helm.InstallNamespaced(NginxRepoPath, NginxNamespace, true)
		if err != nil {
			return err
		}

		glg.Debugf("created nginx release on a DO k8s %v in ns %v", rls.Name, rls.Namespace)
	case "minikube":
		tmpDir := "/tmp"
		minikubeCommand := fmt.Sprintf("minikube addons enable ingress")
		err = util.ExecCommand(minikubeCommand, tmpDir)
		if err != nil {
			glg.Error(err)
		}
	default:
		rls, err := o.Helm.InstallNamespaced(NginxRepoPath, NginxNamespace, true)
		if err != nil {
			return err
		}

		glg.Debugf("created nginx release %v in ns %v", rls.Name, rls.Namespace)
	}

	//0. install kubemq
	rls, err := o.Helm.Install(o.QueuePath)
	if err != nil {
		return err
	}

	glg.Debugf("created queue release %v in ns %v", rls.Name, rls.Namespace)

	return nil
}
