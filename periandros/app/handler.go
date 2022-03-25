package app

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"time"
)

type PeriandrosHandler struct {
	Config   *configs.PeriandrosConfig
	Duration time.Duration
	Timeout  time.Duration
}

func (p *PeriandrosHandler) CreateUser() (bool, error) {
	healthy := p.CheckSolonHealth()
	if !healthy {
		return false, fmt.Errorf("solon not available cannot create user")
	}

	solonResponse, err := p.Config.HttpClients.Solon().Register(p.Config.SolonCreationRequest)
	if err != nil {
		return false, err
	}

	return solonResponse.Created, nil
}

func (p *PeriandrosHandler) CheckSolonHealth() bool {
	healthy := false

	ticker := time.NewTicker(p.Duration)
	timeout := time.After(p.Timeout)

	for {
		select {
		case t := <-ticker.C:
			glg.Infof("tick: %s", t)
			response, err := p.Config.HttpClients.Solon().Health()
			if err != nil {
				glg.Errorf("Error getting response: %s", err)
				continue
			}

			healthy = response.Healthy
			if !healthy {
				continue
			}
			ticker.Stop()

		case <-timeout:
			ticker.Stop()
		}
		break
	}

	return healthy
}
