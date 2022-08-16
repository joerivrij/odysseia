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
	client kubernetes.Interface
}

func NewWorkloadClient(kube kubernetes.Interface) (*WorkloadImpl, error) {
	return &WorkloadImpl{client: kube}, nil
}

func (w *WorkloadImpl) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	deploy, err := w.client.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})

	return deploy, err
}

func (w *WorkloadImpl) ListDeployments(namespace string) (*appsv1.DeploymentList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	deployments, err := w.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

func (w *WorkloadImpl) UpdateDeploymentViaAnnotation(namespace, name string, annotation map[string]string) (*appsv1.Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	deployment, err := w.GetDeployment(namespace, name)
	if err != nil {
		return nil, err
	}

	for key, value := range annotation {
		deployment.Spec.Template.Annotations[key] = value
	}

	return w.client.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
}

func (w *WorkloadImpl) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	return w.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
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

func (w *WorkloadImpl) CreateJob(namespace string, spec *batchv1.Job) (*batchv1.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	job, err := w.client.BatchV1().Jobs(namespace).Create(ctx, spec, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return job, nil
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

func (w *WorkloadImpl) ListJobs(namespace string) (*batchv1.JobList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	jobs, err := w.client.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// List lists all pods within your cluster
func (w *WorkloadImpl) List(namespace string) (*corev1.PodList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	pods, err := w.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (w *WorkloadImpl) DeletePod(namespace, podName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := w.client.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{})

	return err
}

func (w *WorkloadImpl) CreatePod(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	pod, err := w.client.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})

	return pod, err
}

func (w *WorkloadImpl) CreatePodSpec(namespace, name, podImage string, command []string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            name,
					Image:           podImage,
					ImagePullPolicy: corev1.PullAlways,
					Command:         command,
				},
			},
		},
	}

	return pod
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
