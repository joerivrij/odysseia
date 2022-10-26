package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"time"
)

type ThrasyboulosHandler struct {
	Config       *configs.ThrasyboulosConfig
	Duration     time.Duration
	TimeFinished int64
}

func (t *ThrasyboulosHandler) WaitForJobsToFinish(c chan bool) {
	start := time.Now()
	ticker := time.NewTicker(t.Duration)
	defer ticker.Stop()

	for ts := range ticker.C {
		if ts.Sub(start).Milliseconds() >= t.TimeFinished {
			c <- false
		}

		glg.Infof("job: %s is still running", t.Config.Job)

		job, err := t.Config.Kube.Workload().GetJob(t.Config.Namespace, t.Config.Job)
		if err != nil {
			glg.Error(err)
		}

		conditionFound := false
		if job.Status.Active == 0 {
			for _, condition := range job.Status.Conditions {
				if condition.Type == "Complete" {
					conditionFound = true
				}
			}
		}

		if conditionFound {
			c <- true
		}
	}

}
