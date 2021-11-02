package kubernetes

import (
	"archive/tar"
	"bytes"
	"context"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"io"
	"io/ioutil"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type KubeClientImpl struct {
	kubeClient *Kube
}

func (k *KubeClientImpl) GetDeploymentStatus(namespace string) (bool, error) {
	finished := false

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	deployments, err := k.kubeClient.GetK8sClientSet().AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	var deploymentStatus []string

	for _, deployment := range deployments.Items {
		if *deployment.Spec.Replicas == deployment.Status.ReadyReplicas {
			deploymentStatus = append(deploymentStatus, deployment.Name)
		}
	}

	if len(deploymentStatus) == len(deployments.Items) {
		finished = true
	}

	return finished, nil
}

func (k *KubeClientImpl) GetStatefulSets(namespace string) (*v12.StatefulSetList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	sets, err := k.kubeClient.GetK8sClientSet().AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return sets, nil
}


func (k *KubeClientImpl) GetPodsBySelector(namespace, selector string) (*v1.PodList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	pods, err := k.kubeClient.GetK8sClientSet().CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		TypeMeta:            metav1.TypeMeta{},
		LabelSelector:       selector,
		FieldSelector:       "",
		Watch:               false,
		AllowWatchBookmarks: false,
		ResourceVersion:     "",
		TimeoutSeconds:      nil,
		Limit:               0,
		Continue:            "",
	})
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (k *KubeClientImpl) ExecNamedPod(namespace, podName string, command []string) (string, error) {
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	k.kubeClient.GetConfig()

	restCfg, err := kubeCfg.ClientConfig()
	if err != nil {
		return "", err
	}

	commandBuffer := &bytes.Buffer{}
	errBuffer := &bytes.Buffer{}

	req := k.kubeClient.GetK8sClientSet().CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(namespace).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: command,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", req.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: commandBuffer,
		Stderr: errBuffer,
	})

	if err != nil {
		return errBuffer.String(), err
	}

	return commandBuffer.String(), nil
}

func (k *KubeClientImpl) GetSecrets(namespace string) (*v1.SecretList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	secrets, err := k.kubeClient.GetK8sClientSet().CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
		TypeMeta:            metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		LabelSelector:       "",
		FieldSelector:       "",
		Watch:               false,
		AllowWatchBookmarks: false,
		ResourceVersion:     "",
		TimeoutSeconds:      nil,
		Limit:               0,
		Continue:            "",
	})
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (k *KubeClientImpl) CreateSecret(namespace, secretName string, data map[string][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	secret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Immutable:  nil,
		Data:       data,
		StringData: nil,
		Type:       v1.SecretTypeOpaque,
	}

	_, err := k.kubeClient.GetK8sClientSet().CoreV1().Secrets(namespace).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (k *KubeClientImpl) GetServiceAccounts(namespace string) (*v1.ServiceAccountList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	serviceAccounts, err := k.kubeClient.GetK8sClientSet().CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{
		TypeMeta:            metav1.TypeMeta{},
		LabelSelector:       "",
		FieldSelector:       "",
		Watch:               false,
		AllowWatchBookmarks: false,
		ResourceVersion:     "",
		TimeoutSeconds:      nil,
		Limit:               0,
		Continue:            "",
	})
	if err != nil {
		return nil, err
	}

	return serviceAccounts, nil
}

func (k *KubeClientImpl) GetHostServer() (string, error) {
	var server string
	kubeConfig := k.kubeClient.GetConfig()
	config, err := models.UnmarshalKubeConfig(kubeConfig)
	if err != nil {
		return "", err
	}

	currentCtx := config.CurrentContext

	for _, cluster := range config.Clusters {
		if cluster.Name == currentCtx {
			server = cluster.Cluster.Server
		}
	}

	return server, nil
}

func (k *KubeClientImpl) GetHostCaCert() ([]byte, error) {
	kubeConfig := k.kubeClient.GetConfig()
	config, err := models.UnmarshalKubeConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	currentCtx := config.CurrentContext

	for _, cluster := range config.Clusters {
		if cluster.Name == currentCtx {
			if cluster.Cluster.CertificateAuthorityData == "" {
				filePath := cluster.Cluster.CertificateAuthority
				content, err := ioutil.ReadFile(filePath)
				if err != nil {
					return nil, err
				}
				return content, nil
			} else {
				return []byte(cluster.Cluster.CertificateAuthorityData), nil
			}
		}
	}

	return nil, nil
}

func (k *KubeClientImpl) CopyFileToPod(namespace, podName, destPath, srcPath string) (string, error) {
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	k.kubeClient.GetConfig()

	restCfg, err := kubeCfg.ClientConfig()
	if err != nil {
		return "", err
	}

	commandBuffer := &bytes.Buffer{}
	errBuffer := &bytes.Buffer{}

	reader, writer := io.Pipe()
	if destPath != "/" && strings.HasSuffix(string(destPath[len(destPath)-1]), "/") {
		destPath = destPath[:len(destPath)-1]
	}

	go func() {
		defer writer.Close()
		err := makeTar(srcPath, destPath, writer)
		glg.Debug(err)
	}()
	var cmdArr []string

	cmdArr = []string{"tar", "-xf", "-"}
	destDir := path.Dir(destPath)
	if len(destDir) > 0 {
		cmdArr = append(cmdArr, "-C", destDir)
	}

	req := k.kubeClient.GetK8sClientSet().CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(namespace).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: cmdArr,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", req.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  reader,
		Stdout: commandBuffer,
		Stderr: errBuffer,
		Tty:    false,
	})

	if err != nil {
		return errBuffer.String(), err
	}

	return commandBuffer.String(), nil
}

func makeTar(srcPath, destPath string, writer io.Writer) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fpath := range matchedPaths {
		stat, err := os.Lstat(fpath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(fpath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fpath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fpath)
			target, err := os.Readlink(fpath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fpath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
		}
	}
	return nil
}