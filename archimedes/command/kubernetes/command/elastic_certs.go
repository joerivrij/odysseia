package command

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/command"
	"github.com/odysseia/plato/kubernetes"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
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
				glg.Debugf("defaulting to %s", command.DefaultNamespace)
				namespace = command.DefaultNamespace
			}

			if filePath == "" {
				glg.Debugf("defaulting to %s", command.DefaultKubeConfig)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					glg.Error(err)
				}

				filePath = filepath.Join(homeDir, command.DefaultKubeConfig)
			}

			cfg, err := ioutil.ReadFile(filePath)
			if err != nil {
				glg.Error("error getting kubeconfig")
			}

			kubeManager, err := kubernetes.NewKubeClient(cfg, namespace)
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

func createElasticCerts(namespace string, kube kubernetes.KubeClient) {
	glg.Error("not implemented yet")
}

func CreateElasticP12(kube kubernetes.KubeClient, ns, filePath string) (string, error) {
	fileName := "elastic-certificates.p12"
	dns := "elastic-master,localhost,odysseia-greek.internal"
	outPath := fmt.Sprintf("/tmp/%s", fileName)
	commandLine := fmt.Sprintf("elasticsearch-certutil ca --out /tmp/elastic-stack-ca.p12 --pass '' && elasticsearch-certutil cert --name elastic-master --dns %s --ca /tmp/elastic-stack-ca.p12 --pass '' --ca-pass '' --out %s && sleep 50", dns, outPath)
	podName := "elastic"
	podCommand := []string{
		"/bin/sh",
		"-c",
		commandLine,
	}

	podExists, _ := kube.Workload().GetPodByName(ns, podName)
	if podExists != nil {
		kube.Workload().DeletePod(ns, podName)
	}

	podSpec := kube.Workload().CreatePodSpec(ns, podName, "docker.elastic.co/elasticsearch/elasticsearch:7.15.0", podCommand)

	pod, err := kube.Workload().CreatePod(ns, podSpec)
	if err != nil {
		return "", err
	}

	// Wait for 60sec for pod to become ready.
	ticker := time.NewTicker(1 * time.Second)
	timeout := time.After(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			po, err := kube.Workload().GetPodByName(ns, pod.Name)
			if err != nil {
				continue
			}

			if po.Status.Phase != "Running" {
				continue
			}

			ticker.Stop()

		case <-timeout:
			glg.Error("timed out")
			ticker.Stop()
		}
		break
	}

	glg.Infof("pod %s is healthy, extracting file", pod.Name)

	defer func() {
		err = kube.Workload().DeletePod(ns, pod.Name)
		if err != nil {
			glg.Error(err)
		}
	}()

	glg.Info("waiting for process to finish in pod...")
	sleepBeforeCopy := 10 * time.Second
	time.Sleep(sleepBeforeCopy)
	glg.Info("process done")

	err = kube.Util().CopyFileFromPod(outPath, filePath, ns, podName)
	if err != nil {
		return "", err
	}

	srcPath := filepath.Join(filePath, fileName)
	return srcPath, nil
}

func GenerateCrtFromPem(pemBytes []byte) ([]byte, error) {
	caPEM := new(bytes.Buffer)
	err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: pemBytes,
	})

	if err != nil {
		return nil, err
	}

	return caPEM.Bytes(), nil
}
