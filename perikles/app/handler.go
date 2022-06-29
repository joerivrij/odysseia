package app

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type PeriklesHandler struct {
	Config *configs.PeriklesConfig
}

// pingPong pongs the ping
func (p *PeriklesHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// validate that new deployments have the correct secret attached to them
func (p *PeriklesHandler) validate(w http.ResponseWriter, req *http.Request) {
	var body []byte
	if req.Body != nil {
		if data, err := ioutil.ReadAll(req.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		glg.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	arRequest := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &arRequest); err != nil {
		glg.Error("incorrect body")
		http.Error(w, "incorrect body", http.StatusBadRequest)
	}

	raw := arRequest.Request.Object.Raw
	glg.Debug(string(raw))
	deploy := v1.Deployment{}
	if err := json.Unmarshal(raw, &deploy); err != nil {
		glg.Error("error deserializing deployment")
		return
	}

	go func() {
		err := p.checkForAnnotations(deploy)
		if err != nil {
			glg.Error(err)
		}
	}()

	review := v1beta1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       arRequest.Kind,
			APIVersion: arRequest.APIVersion,
		},
		Response: &v1beta1.AdmissionResponse{
			UID:     arRequest.Request.UID,
			Allowed: true,
		},
	}

	middleware.ResponseWithCustomCode(w, 200, review)
}
