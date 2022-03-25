package app

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"net/http"
	"strconv"
	"time"
)

type PtolemaiosHandler struct {
	Config   *configs.PtolemaiosConfig
	Duration time.Duration
}

// PingPong pongs the ping
func (p *PtolemaiosHandler) PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (p *PtolemaiosHandler) Health(w http.ResponseWriter, r *http.Request) {
	vaultHealth, _ := p.Config.Vault.Health()
	glg.Debugf("%s : %s", "vault healthy", strconv.FormatBool(vaultHealth))

	healthy := helpers.GetHealthWithVault(vaultHealth)
	middleware.ResponseWithJson(w, healthy)
}

func (p *PtolemaiosHandler) GetSecretFromVault(w http.ResponseWriter, req *http.Request) {
	oneTimeToken, err := p.getOneTimeToken()
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "getToken",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	p.Config.Vault.SetOnetimeToken(oneTimeToken)
	secret, err := p.Config.Vault.GetSecret(p.Config.PodName)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "getSecret",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var elasticModel models.ElasticConfigVault
	for key, value := range secret.Data {
		if key == "data" {
			j, _ := json.Marshal(value)
			elasticModel, _ = models.UnmarshalSecretData(j)
		}
	}

	middleware.ResponseWithJson(w, elasticModel)

	return
}

func (p *PtolemaiosHandler) getOneTimeToken() (string, error) {
	response, err := p.Config.HttpClients.Solon().OneTimeToken()
	if err != nil {
		return "", err
	}

	glg.Debugf("found token: %s", response.Token)
	return response.Token, nil
}

func (p *PtolemaiosHandler) CheckForJobExit(exitChannel chan bool) {
	var counter int
	for {
		counter++
		glg.Debugf("run number: %d", counter)
		time.Sleep(p.Duration)
		pod, err := p.Config.Kube.Workload().GetPodByName(p.Config.Namespace, p.Config.FullPodName)
		if err != nil {
			glg.Errorf("error getting kube response %s", err)
			continue
		}

		for _, container := range pod.Status.ContainerStatuses {
			if container.Name == p.Config.PodName {
				glg.Debug(container.Name)
				if container.State.Terminated == nil {
					glg.Debugf("%s not done yet", container.Name)
					continue
				}
				if container.State.Terminated.ExitCode == 0 {
					exitChannel <- true
				}
			}
		}
	}
}
