package models

import "encoding/json"

func UnmarshalSolonCreationRequest(data []byte) (SolonCreationRequest, error) {
	var r SolonCreationRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SolonCreationRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SolonCreationRequest struct {
	Role     string   `json:"roles"`
	Access   []string `json:"access"`
	PodName  string   `json:"podName"`
	Username string   `json:"username"`
}

type SolonResponse struct {
	Created bool `json:"created"`
}

func (r *SolonResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *TokenResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TokenResponse struct {
	Token string `json:"token"`
}
