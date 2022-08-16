package kubernetes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPodsClient(t *testing.T) {
	ns := "odysseia"
	expectedName := "testpod"
	access := "dictionary"
	role := "api"

	t.Run("Create", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		podObject := CreatePodObject(expectedName, ns, access, role)

		pod, err := testClient.Workload().CreatePod(ns, podObject)
		assert.Nil(t, err)
		assert.Equal(t, expectedName, pod.Name)
		assert.Equal(t, ns, pod.Namespace)
	})

	t.Run("CreatingTwiceReturnsError", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		podObject := CreatePodObject(expectedName, ns, access, role)

		_, err = testClient.Workload().CreatePod(ns, podObject)
		assert.Nil(t, err)
		_, err = testClient.Workload().CreatePod(ns, podObject)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), expectedName)
	})

	t.Run("Get", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		// we need to create a pod first
		err = CreatePodForTest(expectedName, ns, access, role, testClient)
		assert.Nil(t, err)

		pod, err := testClient.Workload().GetPodByName(ns, expectedName)
		assert.Nil(t, err)
		assert.Equal(t, expectedName, pod.Name)
		assert.Equal(t, ns, pod.Namespace)
	})

	t.Run("List", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)
		expectedNamePodOne := fmt.Sprintf("%sone", expectedName)
		expectedNamePodTwo := fmt.Sprintf("%stwo", expectedName)

		// we need to create a pod first
		err = CreatePodForTest(expectedNamePodOne, ns, access, role, testClient)
		assert.Nil(t, err)
		err = CreatePodForTest(expectedNamePodTwo, ns, access, role, testClient)
		assert.Nil(t, err)

		pods, err := testClient.Workload().List(ns)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(pods.Items))
		for _, pod := range pods.Items {
			assert.Equal(t, pod.Namespace, ns)
			assert.Contains(t, pod.Name, expectedName)

		}
	})
}

func TestJobCLient(t *testing.T) {
	ns := "odysseia"
	expectedName := "testpod"
	completed := true

	t.Run("Create", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		jobObject := CreateJobObject(expectedName, ns, completed)

		job, err := testClient.Workload().CreateJob(ns, jobObject)
		assert.Nil(t, err)
		assert.Equal(t, expectedName, job.Name)
		assert.Equal(t, ns, job.Namespace)
	})

	t.Run("Get", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		jobObject := CreateJobObject(expectedName, ns, completed)

		_, err = testClient.Workload().CreateJob(ns, jobObject)
		assert.Nil(t, err)

		job, err := testClient.Workload().GetJob(ns, expectedName)
		assert.Nil(t, err)
		assert.Equal(t, expectedName, job.Name)
		assert.Equal(t, ns, job.Namespace)
	})

	t.Run("List", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		jobObject := CreateJobObject(expectedName, ns, completed)

		_, err = testClient.Workload().CreateJob(ns, jobObject)
		assert.Nil(t, err)

		jobList, err := testClient.Workload().ListJobs(ns)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(jobList.Items))
		for _, job := range jobList.Items {
			assert.Equal(t, expectedName, job.Name)
			assert.Equal(t, ns, job.Namespace)
		}
	})
}

func TestDeploymentClient(t *testing.T) {
	ns := "odysseia"
	expectedName := "testpod"

	t.Run("Create", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		deployObject := CreateDeploymentObject(expectedName, ns)

		sut, err := testClient.Workload().CreateDeployment(ns, deployObject)
		assert.Nil(t, err)
		assert.Equal(t, expectedName, sut.Name)
		assert.Equal(t, ns, sut.Namespace)
	})

	t.Run("Get", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		deployObject := CreateDeploymentObject(expectedName, ns)

		_, err = testClient.Workload().CreateDeployment(ns, deployObject)
		assert.Nil(t, err)

		sut, err := testClient.Workload().GetDeployment(ns, expectedName)
		assert.Nil(t, err)
		assert.Equal(t, expectedName, sut.Name)
		assert.Equal(t, ns, sut.Namespace)
	})

	t.Run("List", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		deployObject := CreateDeploymentObject(expectedName, ns)

		_, err = testClient.Workload().CreateDeployment(ns, deployObject)
		assert.Nil(t, err)

		sut, err := testClient.Workload().ListDeployments(ns)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(sut.Items))
		for _, deploy := range sut.Items {
			assert.Equal(t, expectedName, deploy.Name)
			assert.Equal(t, ns, deploy.Namespace)
		}
	})

	t.Run("CreateAnnotatedObject", func(t *testing.T) {
		testClient, err := FakeKubeClient(ns)
		assert.Nil(t, err)

		annotations := map[string]string{
			"testkey": "testvalue",
		}
		deployObject := CreateAnnotatedDeploymentObject(expectedName, ns, annotations)

		sut, err := testClient.Workload().CreateDeployment(ns, deployObject)
		assert.Nil(t, err)
		assert.Equal(t, expectedName, sut.Name)
		assert.Equal(t, ns, sut.Namespace)
		assert.Equal(t, annotations, sut.Spec.Template.Annotations)
	})
}
