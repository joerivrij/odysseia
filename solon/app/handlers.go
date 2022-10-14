package app

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/generator"
	"github.com/odysseia-greek/plato/helpers"
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia/aristoteles/configs"
	"net/http"
	"strconv"
	"strings"
)

type SolonHandler struct {
	Config *configs.SolonConfig
}

// PingPong pongs the ping
func (s *SolonHandler) PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (s *SolonHandler) Health(w http.ResponseWriter, r *http.Request) {
	vaultHealth, _ := s.Config.Vault.Health()
	glg.Debugf("%s : %s", "vault healthy", strconv.FormatBool(vaultHealth))

	healthy := helpers.GetHealthWithVault(vaultHealth)
	middleware.ResponseWithJson(w, healthy)
}

func (s *SolonHandler) CreateOneTimeToken(w http.ResponseWriter, req *http.Request) {
	//validate podname as registered?
	policy := []string{"ptolemaios"}
	token, err := s.Config.Vault.CreateOneTimeToken(policy)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "getting token",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

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

	glg.Debug(creationRequest)
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

	glg.Debug("checking pod for correct label")
	//check if pod has the correct labels
	pod, err := s.Config.Kube.Workload().GetPodByName(s.Config.Namespace, creationRequest.PodName)
	var validAccess bool
	var validRole bool
	for key, value := range pod.Annotations {
		if key == s.Config.AccessAnnotation {
			splittedValues := strings.Split(value, ";")
			for _, a := range creationRequest.Access {
				contains := sliceContains(splittedValues, a)
				if !contains {
					break
				}
				glg.Debugf("requested %s matched in annotations %s", a, splittedValues)
				validAccess = contains
			}

		} else if key == s.Config.RoleAnnotation {
			if value == creationRequest.Role {
				glg.Debugf("requested %s matched annotation %s", creationRequest.Role, value)
				validRole = true
			}
		} else {
			continue
		}
	}

	if !validAccess || !validRole {
		glg.Debugf("annotations found on pod %s did not match requested", creationRequest.PodName)
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "annotations",
					Message: fmt.Sprintf("annotations requested and found on pod %s did not match", creationRequest.PodName),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	glg.Debugf("annotations found on pod %s matched requested", creationRequest.PodName)

	var roleNames []string
	for _, a := range creationRequest.Access {
		roleName := fmt.Sprintf("%s_%s", a, creationRequest.Role)
		glg.Debugf("adding role named: %s to user", roleName)
		roleNames = append(roleNames, roleName)
	}

	putUser := elastic.CreateUserRequest{
		Password: password,
		Roles:    roleNames,
		FullName: creationRequest.Username,
		Email:    fmt.Sprintf("%s@odysseia-greek.com", creationRequest.Username),
		Metadata: &elastic.Metadata{Version: 1},
	}

	var response models.SolonResponse
	userCreated, err := s.Config.Elastic.Access().CreateUser(creationRequest.Username, putUser)
	if err != nil {
		glg.Error(err)
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "createUser",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	createRequest := models.CreateSecretRequest{
		Data: models.ElasticConfigVault{
			Username:    creationRequest.Username,
			Password:    password,
			ElasticCERT: string(s.Config.ElasticCert),
		},
	}

	payload, _ := createRequest.Marshal()

	secretCreated, err := s.Config.Vault.CreateNewSecret(creationRequest.Username, payload)
	if err != nil {
		glg.Error(err)
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "createSecret",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	glg.Debugf("secret created in vault %t", secretCreated)

	response.Created = userCreated

	middleware.ResponseWithJson(w, response)
	return
}

func sliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
