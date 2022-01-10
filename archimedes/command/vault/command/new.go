package command

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
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

			config, err := aristoteles.NewConfig()
			if err != nil {
				glg.Error(err)
				os.Exit(1)
			}

			kube, err := config.GetKubeClient(filePath, namespace)
			if err != nil {
				glg.Fatal("error creating kubeclient")
			}

			glg.Info("1. vault init started")
			initVault(namespace, kube)
			glg.Info("1. vault init completed")
			glg.Info("2. vault unseal started")
			unsealVault("", namespace, kube)
			glg.Info("2. vault unseal completed")
			glg.Info("3. adding admin")
			createPolicy(defaultAdminPolicyName, namespace, kube)
			glg.Info("3. finished adding admin")
			glg.Info("4. adding user")
			createPolicy(defaultUserPolicyName, namespace, kube)
			glg.Info("4. finished adding user")
			glg.Info("5. adding kuberentes as auth method")
			enableKubernetesAsAuth(namespace, defaultAdminPolicyName, kube)
			glg.Info("5. finished adding kuberentes as auth method")

		},
	}

	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "kubeconfig filepath defaults to ~/.kube/config")

	return cmd
}
