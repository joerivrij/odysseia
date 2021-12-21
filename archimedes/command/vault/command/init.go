package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/util"
	"github.com/odysseia/plato/configuration"
	"github.com/odysseia/plato/kubernetes"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Init() *cobra.Command {
	var (
		namespace string
		filePath  string
	)
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Inits your vault",
		Long: `Allows you to init the vault, it takes
- Namespace
- Filepath`,
		Run: func(cmd *cobra.Command, args []string) {

			if namespace == "" {
				glg.Debugf("defaulting to %s", defaultNamespace)
				namespace = defaultNamespace
			}

			if filePath == "" {
				glg.Debugf("defaulting to %s", defaultKubeConfig)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					glg.Error(err)
				}

				filePath = filepath.Join(homeDir, defaultKubeConfig)
			}

			config, err := configuration.NewConfig()
			if err != nil {
				glg.Error(err)
				os.Exit(1)
			}

			kube, err := config.GetKubeClient(filePath, namespace)
			if err != nil {
				glg.Fatal("error creating kubeclient")
			}

			glg.Info("is it secret? Is it safe? Well no longer!")
			glg.Debug("unsealing kube vault")
			initVault(namespace, kube)
		},
	}

	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "kubeconfig filepath defaults to ~/.kube/config")

	return cmd
}

func initVault(namespace string, kube *kubernetes.Kube) {
	vaultSelector := "app.kubernetes.io/name=vault"
	var podName string

	pods, err := kube.Workload().GetPodsBySelector(namespace, vaultSelector)
	if err != nil {
		glg.Error(err)
	}
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "vault") {
			if pod.Status.Phase == "Running" {
				glg.Debugf(fmt.Sprintf("%s is running in release %s", pods.Items[0].Name, namespace))
				podName = pod.Name
				break
			}
		}
	}

	command := []string{"vault", "operator", "init", "-key-shares=1", "-key-threshold=1", "-format=json"}

	vaultInit, err := kube.Workload().ExecNamedPod(namespace, podName, command)
	if err != nil {
		glg.Error(err)
		return
	}

	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == "odysseia" {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}

	fileName := fmt.Sprintf("cluster-keys-%s.json", namespace)
	clusterKeys := filepath.Join(l, "solon", "vault_config", fileName)

	util.WriteFile([]byte(vaultInit), clusterKeys)

	glg.Debugf("wrote data: %s to dest: %s", vaultInit, clusterKeys)
}
