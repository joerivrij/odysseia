package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/vault"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

const defaultNamespace = "odysseia"
const defaultKubeConfig = "/.kube/config"

func Unseal() *cobra.Command {
	var (
		key       string
		namespace string
		filePath  string
	)
	cmd := &cobra.Command{
		Use:   "unseal",
		Short: "Unseal your vault",
		Long: `Allows you unseal the vault, it takes
- Key`,
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

			baseConfig := configs.ArchimedesConfig{}
			unparsedConfig, err := aristoteles.NewConfig(baseConfig)
			if err != nil {
				glg.Error(err)
				glg.Fatal("death has found me")
			}
			archimedesConfig, ok := unparsedConfig.(*configs.ArchimedesConfig)
			if !ok {
				glg.Fatal("could not parse config")
			}

			glg.Info("is it secret? Is it safe? Well no longer!")
			glg.Debug("unsealing kube vault")
			UnsealVault(key, namespace, archimedesConfig.Kube)
		},
	}

	cmd.PersistentFlags().StringVarP(&key, "key", "k", "", "unseal key, if not set cluster-keys will be used")
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "kubeconfig filepath defaults to ~/.kube/config")

	return cmd
}

func UnsealVault(key, namespace string, kube kubernetes.KubeClient) {
	if key == "" {
		glg.Info("key was not given, trying to get key from cluster-keys.json")
		clusterKeys, err := vault.GetClusterKeys()
		if err != nil {
			glg.Fatal("could not get cluster keys")
		}
		key = clusterKeys.VaultUnsealKey
		glg.Info("key found")
	}

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

	command := []string{"vault", "operator", "unseal", key}

	vaultUnsealed, err := kube.Workload().ExecNamedPod(namespace, podName, command)
	if err != nil {
		glg.Error(err)
	}

	glg.Info(vaultUnsealed)
}
