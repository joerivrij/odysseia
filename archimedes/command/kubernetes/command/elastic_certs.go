package command

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/kubernetes"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func CreateElasticCerts() *cobra.Command {
	var (
		namespace string
		filePath  string
	)
	cmd := &cobra.Command{
		Use:   "elastic_certs",
		Short: "Create the certs needed for a ssl setup in Elastic",
		Long: `Allows you to create a elastic cert without using the elastic utils
- Namespace
- Filepth`,
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

			kubeManager, err := kubernetes.NewKubeClient(filePath, namespace)
			if err != nil {
				glg.Fatal("error creating kubeclient")
			}

			glg.Info("is it secret? Is it safe? Well no longer!")
			glg.Debug("creating elastic certs")

			createElasticCerts(namespace, kubeManager)
		},
	}

	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "kubeconfig filepath defaults to ~/.kube/config")

	return cmd
}

func createElasticCerts(namespace string, kube *kubernetes.Kube) {
	glg.Error("not implemented yet")
}
