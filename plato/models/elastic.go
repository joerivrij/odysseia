package models

import "encoding/json"

func UnmarshalElasticResponse(data []byte) (ElasticResponse, error) {
	var r ElasticResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ElasticResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ElasticResponse struct {
	ScrollId string `json:"_scroll_id,omitempty"`
	Took     int64  `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
}

type Hits struct {
	Total    Total   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Hit struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

type Shards struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

func UnmarshalCreateRoleRequest(data []byte) (CreateRoleRequest, error) {
	var r CreateRoleRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CreateRoleRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CreateRoleRequest struct {
	Cluster      []string      `json:"cluster"`
	Indices      []Index       `json:"indices"`
	Applications []Application `json:"applications"`
	RunAs        []string      `json:"run_as,omitempty"`
	Metadata     Metadata      `json:"metadata,omitempty"`
}

type Application struct {
	Application string   `json:"application"`
	Privileges  []string `json:"privileges"`
	Resources   []string `json:"resources"`
}

type Index struct {
	Names         []string       `json:"names"`
	Privileges    []string       `json:"privileges"`
	FieldSecurity *FieldSecurity `json:"field_security,omitempty"`
	Query         string         `json:"query,omitempty"`
}

type FieldSecurity struct {
	Grant []string `json:"grant"`
}

type Metadata struct {
	Version int64 `json:"version,omitempty"`
}

func UnmarshalCreateUserRequest(data []byte) (CreateUserRequest, error) {
	var r CreateUserRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CreateUserRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CreateUserRequest struct {
	Password string    `json:"password"`
	Roles    []string  `json:"roles"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Metadata *Metadata `json:"metadata"`
}
