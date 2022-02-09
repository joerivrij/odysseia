package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/kubernetes"
	"strings"
)

func enableSecrets(namespace, configName, rootToken string, kube kubernetes.KubeClient) {
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

	loginCommand := []string{"vault", "login", rootToken}

	login, err := kube.Workload().ExecNamedPod(namespace, podName, loginCommand)
	if err != nil {
		glg.Error(err)
	}

	glg.Info(login)

	path := fmt.Sprintf("-path=%s", configName)
	command := []string{"vault", "secrets", "enable", path, "kv"}

	secretEnabled, err := kube.Workload().ExecNamedPod(namespace, podName, command)
	if err != nil {
		glg.Error(err)
	}

	glg.Info(secretEnabled)
}
