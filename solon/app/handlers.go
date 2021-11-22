package app

import (
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/generator"
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

	healthy := helpers.GetHealthWithVault(true)
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
	var creationRequest models.SolonCreationRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&creationRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "decoding",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	password, err := generator.RandomPassword(18)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "passwordgenerator",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	//check if pod has the correct labels
	namespace := "odysseia"
	pod, err := s.Config.Kube.GetPodByName(namespace, creationRequest.PodName)
	glg.Debug(pod.Annotations)


	var roleNames []string
	for _, a := range creationRequest.Access {
		roleName := fmt.Sprintf("%s_%s", a, creationRequest.Role)
		glg.Debugf("adding new role named: %s to elastic", roleName)
		roleNames = append(roleNames, roleName)
	}

	putUser := models.CreateUserRequest{
		Password: password,
		Roles:    roleNames,
		FullName: creationRequest.PodName,
		Email:    fmt.Sprintf("%s@odysseia-greek.com", creationRequest.PodName),
		Metadata: &models.Metadata{Version: 1},
	}

	var response models.SolonResponse
	userCreated, err := elastic.CreateUser(&s.Config.ElasticClient, creationRequest.PodName, putUser)
	if err != nil {
		glg.Error(err)
	}

	createRequest := models.CreateSecretRequest{
		Data: models.SecretData{
			Username:    creationRequest.PodName,
			Password:    password,
			ElasticCERT: s.Config.ElasticCert,
		},
	}

	payload, _ := createRequest.Marshal()

	secretCreated, _ := s.createNewSecret(creationRequest.PodName, payload)
	glg.Debugf("secret created in vault %t", secretCreated)

	response.Created = userCreated

	middleware.ResponseWithJson(w, response)
	return
}

func (s *SolonHandler)vaultHealth() (bool, error) {
	_, err := s.Config.VaultClient.Logical().Read("/sys/health")
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}


func (s *SolonHandler)createNewSecret(name string, payload []byte) (bool, error) {
	vaultPath := fmt.Sprintf("configs/data/%s", name)

	secret, err := s.Config.VaultClient.Logical().WriteBytes(vaultPath, payload)
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	glg.Debug(secret.Data)

	return true, nil
}