package kubernetes

import (
	"bytes"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/remotecommand"
	"time"
)

type WorkloadImpl struct {
	client *kubernetes.Clientset
}

func NewWorkloadClient(kube *kubernetes.Clientset) (*WorkloadImpl, error) {
	return &WorkloadImpl{client: kube}, nil
}

func (w *WorkloadImpl) GetDeploymentStatus(namespace string) (bool, error) {
	finished := false

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	deployments, err := w.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
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

func (w *WorkloadImpl) GetStatefulSets(namespace string) (*appsv1.StatefulSetList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	sets, err := w.client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return sets, nil
}

func (w *WorkloadImpl) GetJob(namespace, name string) (*batchv1.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	job, err := w.client.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (w *WorkloadImpl) GetPodsBySelector(namespace, selector string) (*corev1.PodList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	pods, err := w.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
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

func (w *WorkloadImpl) GetPodByName(namespace, name string) (*corev1.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	pod, err := w.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (w *WorkloadImpl) ExecNamedPod(namespace, podName string, command []string) (string, error) {
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	restCfg, err := kubeCfg.ClientConfig()
	if err != nil {
		return "", err
	}

	commandBuffer := &bytes.Buffer{}
	errBuffer := &bytes.Buffer{}

	req := w.client.CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(namespace).SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
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

func (w *WorkloadImpl) GetNewLock(lockName, podName, namespace string) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockName,
			Namespace: namespace,
		},
		Client: w.client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podName,
		},
	}
}
