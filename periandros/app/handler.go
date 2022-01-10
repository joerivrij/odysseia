package app

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/models"
	"time"
)

type PeriandrosHandler struct {
	Config *configs.PeriandrosConfig
}

func (p *PeriandrosHandler) CreateUser() (bool, error) {
	path := "solon/v1/register"
	p.Config.SolonService.Path = path

	body, _ := p.Config.SolonCreationRequest.Marshal()

	response, err := helpers.PostRequest(p.Config.SolonService, body)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	var solonResponse models.SolonResponse
	err = json.NewDecoder(response.Body).Decode(&solonResponse)
	if err != nil {
		return false, err
	}

	return solonResponse.Created, nil
}

func (p *PeriandrosHandler) CheckSolonHealth(ticks time.Duration) bool {
	path := "solon/v1/health"
	p.Config.SolonService.Path = path

	healthy := false

	ticker := time.NewTicker(1 * time.Second)
	timeout := time.After(ticks * time.Second)

	for {
		select {
		case t := <-ticker.C:
			glg.Infof("tick: %s", t)
			response, err := helpers.GetRequest(p.Config.SolonService)
			if err != nil {
				glg.Errorf("Error getting response: %s", err)
				continue
			}

			defer response.Body.Close()

			var healthResponse models.Health
			err = json.NewDecoder(response.Body).Decode(&healthResponse)
			if err != nil {
				glg.Errorf("Error getting response: %s", err)
				continue
			}

			healthy = healthResponse.Healthy
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
