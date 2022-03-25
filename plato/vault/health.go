package vault

import (
	"fmt"
	"github.com/kpango/glg"
	"time"
)

func (v *Vault) Health() (bool, error) {
	_, err := v.Connection.Logical().Read("sys/health")
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}

func (v *Vault) CheckHealthyStatus(ticks, tick time.Duration) bool {
	healthy := false

	ticker := time.NewTicker(tick)
	timeout := time.After(ticks)

	for {
		select {
		case t := <-ticker.C:
			glg.Infof("tick: %s", t)
			res, err := v.Health()
			if err != nil {
				glg.Errorf("Error getting response: %s", err)
				continue
			}
			if res {
				healthy = true
				ticker.Stop()
			}

		case <-timeout:
			ticker.Stop()
		}
		break
	}

	return healthy
}
