package command

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/generator"
	"github.com/odysseia/plato/kubernetes"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const defaultKubeConfig = "/.kube/config"
const defaultNamespace = "odysseia"
const defaultSecretName = "elastic-root-secret"

func CreateSecret() *cobra.Command {
	var (
		secretName   string
		namespace    string
		filePath     string
		secretLength int
	)
	cmd := &cobra.Command{
		Use:   "create_secret",
		Short: "Create a secret",
		Long: `Allows you to create a random secret
- SecretName
- Namespace
- FilePath`,
		Run: func(cmd *cobra.Command, args []string) {

			if namespace == "" {
				glg.Debugf("defaulting to %s", defaultNamespace)
				namespace = defaultNamespace
			}

			if secretName == "" {
				glg.Debugf("defaulting to %s", defaultSecretName)
				secretName = defaultSecretName
			}

			if filePath == "" {
				glg.Debugf("defaulting to %s", defaultKubeConfig)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					glg.Error(err)
				}

				filePath = filepath.Join(homeDir, defaultKubeConfig)
			}

			kubeManager, err := kubernetes.NewKubeClient(filePath)
			if err != nil {
				glg.Fatal("error creating kubeclient")
			}

			glg.Info("is it secret? Is it safe? Well no longer!")
			glg.Debug("creating a kube secret")
			createSecret(secretName, namespace, secretLength, kubeManager)
		},
	}

	cmd.PersistentFlags().StringVarP(&secretName, "secret_name", "s", "", "secret name")
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "kubeconfig filepath defaults to ~/.kube/config")
	cmd.PersistentFlags().IntVarP(&secretLength, "length", "l", 24, "lenght of secret to create")

	return cmd
}

func createSecret(secretName, namespace string, secretLength int, kube kubernetes.Client) {
	password, err := generator.RandomPassword(secretLength)
	if err != nil {
		glg.Error(err)
		return
	}

	glg.Debug(password)

	data := make(map[string][]byte)
	data["password"] = []byte(password)
	data["username"] = []byte("elastic")

	err = kube.CreateSecret(namespace, secretName, data)
	if err != nil {
		glg.Error(err)
		return
	}

	glg.Infof("create secret with name %s", secretName)
}
