package app

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"net/http"
	"strconv"
)

type SolonHandler struct {
	Config *SolonConfig
}

// PingPong pongs the ping
func (s *SolonHandler) PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (s *SolonHandler)Health(w http.ResponseWriter, r *http.Request) {
	glg.Debug("connecting to vault to establish health")
	vaultHealth, _ := s.vaultHealth()
	glg.Debugf("%s : %s", "vault healthy", strconv.FormatBool(vaultHealth))

	healthy := helpers.GetHealthWithVault(vaultHealth)
	middleware.ResponseWithJson(w, healthy)
}

func (s *SolonHandler) CreateOneTimeToken(w http.ResponseWriter, req *http.Request) {
	policy := []string{"odysseia"}
	renew := false

	tokenRequest := vault.TokenCreateRequest{
		Policies:        policy,
		TTL:             "5m",
		DisplayName:     "solonCreated",
		NumUses:         1,
		Renewable:       &renew,
	}

	glg.Debug("request created")

	resp, _ := s.Config.VaultClient.Auth().Token().Create(&tokenRequest)
	token := resp.Auth.ClientToken

	glg.Debug(token)

	tokenModel := models.TokenResponse{
		Token: token,
	}

	middleware.ResponseWithJson(w, tokenModel)
}


func (s *SolonHandler) RegisterService(w http.ResponseWriter, req *http.Request) {
	glg.Info(req.RemoteAddr)
	glg.Info(req.Host)
	glg.Info(req.Header)
	glg.Info(req.RequestURI)
	middleware.ResponseWithJson(w, req)
}

func (s *SolonHandler)vaultHealth() (bool, error) {
	_, err := s.Config.VaultClient.Logical().Read("secret/data")
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}