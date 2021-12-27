package app

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/vault"
	"net/http"
	"os"
	"time"
)

type PtolemaiosHandler struct {
	Config *PtolemaiosConfig
}

// PingPong pongs the ping
func (p *PtolemaiosHandler) PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
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

	client, err := vault.CreateVaultClient(p.Config.VaultService, oneTimeToken)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "createVault",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	secret, err := client.GetSecret(p.Config.PodName)
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
	if p.Config.IsPartOfJob {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
}

func (p *PtolemaiosHandler) getOneTimeToken() (string, error) {
	u := p.Config.SolonService
	u.Path = "/solon/v1/token"
	response, err := helpers.GetRequest(u)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	var tokenModel models.TokenResponse
	err = json.NewDecoder(response.Body).Decode(&tokenModel)
	if err != nil {
		return "", err
	}

	glg.Debugf("found token: %s", tokenModel.Token)

	return tokenModel.Token, nil
}
