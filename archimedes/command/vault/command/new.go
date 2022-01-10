package command

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func New() *cobra.Command {
	var (
		namespace string
		filePath  string
	)
	cmd := &cobra.Command{
		Use:   "new",
		Short: "adds the full flow to vault",
		Long:  `inits, unseals vault and adds both policies and auth method`,
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

			glg.Info("1. vault init started")
			initVault(namespace, archimedesConfig.Kube)
			glg.Info("1. vault init completed")
			glg.Info("2. vault unseal started")
			unsealVault("", namespace, archimedesConfig.Kube)
			glg.Info("2. vault unseal completed")
			glg.Info("3. adding admin")
			createPolicy(defaultAdminPolicyName, namespace, archimedesConfig.Kube)
			glg.Info("3. finished adding admin")
			glg.Info("4. adding user")
			createPolicy(defaultUserPolicyName, namespace, archimedesConfig.Kube)
			glg.Info("4. finished adding user")
			glg.Info("5. adding kuberentes as auth method")
			enableKubernetesAsAuth(namespace, defaultAdminPolicyName, archimedesConfig.Kube)
			glg.Info("5. finished adding kuberentes as auth method")

		},
	}

	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "kubeconfig filepath defaults to ~/.kube/config")

	return cmd
}
