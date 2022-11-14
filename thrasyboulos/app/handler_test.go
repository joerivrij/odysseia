package app

import (
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/kubernetes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJobExit(t *testing.T) {
	ns := "odysseia"
	expectedName := "testpod"
	duration := 10 * time.Millisecond
	timeFinished := int64(1000)

	t.Run("JobFinished", func(t *testing.T) {
		testClient, err := kubernetes.FakeKubeClient(ns)
		assert.Nil(t, err)

		jobSpec := kubernetes.CreateJobObject(expectedName, ns, true)
		job, err := testClient.Workload().CreateJob(ns, jobSpec)
		assert.Nil(t, err)
		assert.Equal(t, job.Name, expectedName)

		testConfig := configs.ThrasyboulosConfig{
			Kube:      testClient,
			Job:       expectedName,
			Namespace: ns,
		}

		handler := ThrasyboulosHandler{Config: &testConfig, Duration: duration, TimeFinished: timeFinished}
		jobExit := make(chan bool, 1)
		go handler.WaitForJobsToFinish(jobExit)

		select {

		case <-jobExit:
			exitStatus := <-jobExit
			assert.True(t, exitStatus)
		}
	})

	t.Run("JobNotFinished", func(t *testing.T) {
		testClient, err := kubernetes.FakeKubeClient(ns)
		assert.Nil(t, err)

		jobSpec := kubernetes.CreateJobObject(expectedName, ns, false)
		job, err := testClient.Workload().CreateJob(ns, jobSpec)
		assert.Nil(t, err)
		assert.Equal(t, job.Name, expectedName)

		testConfig := configs.ThrasyboulosConfig{
			Kube:      testClient,
			Job:       expectedName,
			Namespace: ns,
		}

		timeFinished = duration.Milliseconds() * 2

		handler := ThrasyboulosHandler{Config: &testConfig, Duration: duration, TimeFinished: timeFinished}
		jobExit := make(chan bool, 1)
		go handler.WaitForJobsToFinish(jobExit)

		select {

		case <-jobExit:
			exitStatus := <-jobExit
			assert.False(t, exitStatus)
		}
	})
}
