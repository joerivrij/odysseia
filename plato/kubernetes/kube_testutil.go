package kubernetes

import (
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	podKind       string            = "Pod"
	podVersion    string            = "v1"
	podImage      string            = "testimage"
	podPullPolicy corev1.PullPolicy = "Always"
	jobCommand    string            = "ls"
)

func CreatePodObject(name, ns, access, role string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       podKind,
			APIVersion: podVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"app":      name,
				"job-name": name,
			},
			Annotations: map[string]string{
				"odysseia-greek/access": access,
				"odysseia-greek/role":   role,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            name,
					Image:           podImage,
					ImagePullPolicy: podPullPolicy,
				},
			},
		},
		Status: corev1.PodStatus{
			Phase: "Succeeded",
		},
	}

	return pod
}

func CreatePodObjectWithExit(name, ns string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       podKind,
			APIVersion: podVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"app":      name,
				"job-name": name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            name,
					Image:           podImage,
					ImagePullPolicy: podPullPolicy,
				},
			},
		},
		Status: corev1.PodStatus{
			Phase: "Succeeded",
			ContainerStatuses: []corev1.ContainerStatus{{
				Name: name,
				State: corev1.ContainerState{
					Waiting: nil,
					Running: nil,
					Terminated: &corev1.ContainerStateTerminated{
						ExitCode:    0,
						Signal:      0,
						Reason:      "",
						Message:     "",
						StartedAt:   metav1.Time{},
						FinishedAt:  metav1.Time{},
						ContainerID: "",
					},
				},
				LastTerminationState: corev1.ContainerState{},
				Ready:                true,
				RestartCount:         0,
				Image:                "",
				ImageID:              "",
				ContainerID:          "",
				Started:              nil,
			},
			},
		},
	}

	return pod
}

func CreateJobObject(name, ns string, completed bool) *batchv1.Job {
	var conditionType batchv1.JobConditionType
	if completed {
		conditionType = batchv1.JobComplete
	} else {
		conditionType = batchv1.JobSuspended
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    fmt.Sprintf("%s-pod", name),
							Image:   podImage,
							Command: []string{jobCommand},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
		Status: batchv1.JobStatus{
			Conditions: []batchv1.JobCondition{
				{
					Type:               conditionType,
					Status:             "True",
					LastProbeTime:      metav1.Time{},
					LastTransitionTime: metav1.Time{},
					Reason:             "",
					Message:            "",
				},
			},
		},
	}

	return job
}

func CreatePodForTest(name, ns, access, role string, client KubeClient) error {
	pod := CreatePodObject(name, ns, access, role)
	_, err := client.Workload().CreatePod(ns, pod)
	return err
}
