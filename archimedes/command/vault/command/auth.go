package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/util"
	"github.com/odysseia/plato/kubernetes"
	"github.com/odysseia/plato/vault"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

const defaultMethod = "kubernetes"

func Auth() *cobra.Command {
	var (
		method     string
		namespace  string
		filePath   string
		policyName string
	)
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "adds auth methods to vault",
		Long:  `adds auth methods to vault currently only supports kubernetes`,
		Run: func(cmd *cobra.Command, args []string) {
			if namespace == "" {
				glg.Debugf("defaulting to %s", defaultNamespace)
				namespace = defaultNamespace
			}

			if policyName == "" {
				glg.Debugf("defaulting to %s", defaultAdminPolicyName)
				policyName = defaultAdminPolicyName
			}

			if filePath == "" {
				glg.Debugf("defaulting to %s", defaultKubeConfig)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					glg.Error(err)
				}

				filePath = filepath.Join(homeDir, defaultKubeConfig)
			}

			if method == "" {
				glg.Debugf("defaulting to %s", defaultMethod)
				method = defaultMethod
			} else if method != defaultMethod {
				glg.Debugf("only default %s supported", defaultMethod)
				method = defaultMethod
			}

			kubeManager, err := kubernetes.NewKubeClient(filePath)
			if err != nil {
				glg.Fatal("error creating kubeclient")
			}

			glg.Infof("enabling the following auth method %s", method)
			EnableKubernetesAsAuth(namespace, policyName, kubeManager)
		},
	}

	cmd.PersistentFlags().StringVarP(&method, "method", "m", "", "method to enable in vault defaults to kubernetes")
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace defaults to odysseia")
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "kubeconfig filepath defaults to ~/.kube/config")
	cmd.PersistentFlags().StringVarP(&policyName, "policy", "p", "", "policy to add to kubernetes auth default to solon-admin")

	return cmd
}

func EnableKubernetesAsAuth(namespace, policyName string, kube kubernetes.Client) {
	vaultSubString := "vault"
	var vaultSelector string
	var podName string
	var serviceAccountName string

	sets, err := kube.GetStatefulSets(namespace)
	if err != nil {
		glg.Error(err)
	}

	for _, set := range sets.Items {
		if strings.Contains(set.Name, "vault") {
			serviceAccountName = set.Spec.Template.Spec.ServiceAccountName
			for key, element := range set.Spec.Selector.MatchLabels {
				if element == vaultSubString {
					vaultSelector = fmt.Sprintf("%s=%s", key, element)
				}
			}
		}
	}

	pods, err := kube.GetPodsBySelector(namespace, vaultSelector)
	if err != nil {
		glg.Error(err)
	}
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, vaultSubString) {
			if pod.Status.Phase == "Running" {
				glg.Debugf(fmt.Sprintf("%s is running in release %s", pods.Items[0].Name, namespace))
				podName = pod.Name
				break
			}
		}
	}

	glg.Debugf("found vault pod running as %s", podName)

	kubeHost, _ := kube.GetHostServer()

	glg.Debugf("kubehost found: %s", kubeHost)

	secrets, err := kube.GetSecrets(namespace)
	if err != nil {
		glg.Error(err)
	}

	searchString := "vault-token"
	var jwtToken string

	for _, secret := range secrets.Items {
		if strings.Contains(secret.Name, searchString) {
			jwtToken = string(secret.Data["token"])
			break
		}
	}

	glg.Debugf("found token in cluster: %s", jwtToken)

	filePath := "/tmp/ca-cert-vault-archimedes"
	ca, err := kube.GetHostCaCert()
	if err != nil {
		glg.Error(err)
	}

	glg.Debug(ca)

	util.WriteFile(ca, filePath)

	writeResult, err := kube.CopyFileToPod(namespace, podName, filePath, filePath)
	if err != nil {
		glg.Error(err)
	}

	if writeResult == "" {
		glg.Debugf("copied file: %s to pod: %s", filePath, podName)
	}

	glg.Info("gathered all data needed to start and enable kubernetes as auth and adding your current cluster")

	glg.Info("step 0: logging in using roottoken")

	clusterKeys := vault.GetClusterKeys(namespace)
	rootToken := clusterKeys.RootToken
	glg.Info("key found")

	loginCommand := []string{"vault", "login", rootToken}

	login, err := kube.ExecNamedPod(namespace, podName, loginCommand)
	if err != nil {
		glg.Error(err)
	}

	glg.Debug(login)

	glg.Info("step 1: enable kubernetes as auth method")
	command := []string{"vault", "auth", "enable", "kubernetes"}

	enableAuth, err := kube.ExecNamedPod(namespace, podName, command)
	if err != nil {
		glg.Error(err)
	}

	glg.Debug(enableAuth)
	glg.Info("step 1: finished")

	glg.Info("step 2: write config to kubernetes auth method")

	reviewer := fmt.Sprintf("token_review_jwt=%s", jwtToken)
	configHost := fmt.Sprintf("kubernetes_host=%s", kubeHost)
	caCert := fmt.Sprintf("kubernetes_ca_cert=@%s", filePath)

	addConfigCommand := []string{"vault", "write", "auth/kubernetes/config", reviewer, configHost, caCert}

	configAdded, err := kube.ExecNamedPod(namespace, podName, addConfigCommand)
	if err != nil {
		glg.Error(err)
	}

	glg.Debug(configAdded)
	glg.Info("step 2: finished")

	glg.Info("step 3: enabling role on kubernetes auth")
	glg.Info("role will be called Solon")

	//add new serviceaccount: serviceAccountName,access-sa
	boundServiceName := fmt.Sprintf("bound_service_account_names=%s,%s", serviceAccountName, "access-sa")
	boundNamespace := fmt.Sprintf("bound_service_account_namespaces=%s", namespace)
	policies := fmt.Sprintf("policies=%s", policyName)

	roleCommand := []string{"vault", "write", "auth/kubernetes/role/solon", boundServiceName, boundNamespace, policies}

	roleOutput, err := kube.ExecNamedPod(namespace, podName, roleCommand)
	if err != nil {
		glg.Error(err)
	}

	glg.Debug(roleOutput)
	glg.Info("step 3: finished")

	glg.Info("successfully finished creating a kubernetes auth method for vault")
}
