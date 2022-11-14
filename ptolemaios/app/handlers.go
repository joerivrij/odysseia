package app

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia-greek/plato/models"
	pb "github.com/odysseia-greek/plato/proto"
	"net/http"
	"time"
)

type PtolemaiosHandler struct {
	Config   *configs.PtolemaiosConfig
	Duration time.Duration
	pb.UnimplementedPtolemaiosServer
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

	glg.Debug("so far so good")
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
