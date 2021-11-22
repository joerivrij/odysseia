package app

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"net/http"
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
	engine := req.URL.Query().Get("config")
	secretName := req.URL.Query().Get("service")

	if engine == "" || secretName == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "config",
					Message: "cannot be empty",
				},
				{
					Field:   "service",
					Message: "cannot be empty",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	clientToken := "s.tsLsUgLsqQqya4Fhkn9sgeXz"

	vaultURL := fmt.Sprintf("%s/v1/%s/data/%s?metadata=1", p.Config.VaultService, engine, secretName)
	vaultReq, _ := http.NewRequest("GET", vaultURL, nil)
	vaultReq.Header.Add("X-Vault-Token", clientToken)

	secret, err := p.getSecretFromVault(*vaultReq)
	if err != nil {
		middleware.ResponseWithJson(w, err)
		return
	}

	if secret.Data == nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "secret",
					Message: "was not retrievable",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	middleware.ResponseWithJson(w, secret.Data)
}

func (p *PtolemaiosHandler)getSecretFromVault(vaultRequest http.Request) (*vault.Secret, error) {
	client := &http.Client{}
	resp, err := client.Do(&vaultRequest)
	if err != nil {
		return nil, err
	}

	secret, err := vault.ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (p *PtolemaiosHandler)getOneTimeToken() (string, error){

	return "", nil
}

func (p *PtolemaiosHandler)registerServiceAtStartUp()  {

}