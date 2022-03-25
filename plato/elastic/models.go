package elastic

import "encoding/json"

func UnmarshalResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Response struct {
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

func UnmarshalAggregations(data []byte) (Aggregations, error) {
	var r Aggregations
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Aggregations) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Aggregations struct {
	Took         int64              `json:"took"`
	TimedOut     bool               `json:"timed_out"`
	Shards       Shards             `json:"_shards"`
	Hits         Hits               `json:"hits"`
	Aggregations AuthorAggregations `json:"aggregations"`
}

type AuthorAggregations struct {
	AuthorAggregation   Aggregation `json:"authors"`
	BookAggregation     Aggregation `json:"books"`
	CategoryAggregation Aggregation `json:"categories"`
}

type Aggregation struct {
	DocCountErrorUpperBound int64    `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int64    `json:"sum_other_doc_count"`
	Buckets                 []Bucket `json:"buckets"`
}

type Bucket struct {
	Key      interface{} `json:"key"`
	DocCount int64       `json:"doc_count"`
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
	Indices      []Indices     `json:"indices"`
	Applications []Application `json:"applications"`
	RunAs        []string      `json:"run_as,omitempty"`
	Metadata     Metadata      `json:"metadata,omitempty"`
}

type Application struct {
	Application string   `json:"application"`
	Privileges  []string `json:"privileges"`
	Resources   []string `json:"resources"`
}

type Indices struct {
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

type Config struct {
	Service     string `json:"elasticService"`
	Username    string `json:"elasticUsername"`
	Password    string `json:"elasticPassword"`
	ElasticCERT string `json:"elasticCert"`
}

func UnmarshalCreateResult(data []byte) (CreateResult, error) {
	var r CreateResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CreateResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CreateResult struct {
	Index       string `json:"_index"`
	Type        string `json:"_type"`
	ID          string `json:"_id"`
	Version     int64  `json:"_version"`
	Result      string `json:"result"`
	Shards      Shards `json:"_shards"`
	SeqNo       int64  `json:"_seq_no"`
	PrimaryTerm int64  `json:"_primary_term"`
}

func UnmarshalIndexCreateResult(data []byte) (IndexCreateResult, error) {
	var r IndexCreateResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *IndexCreateResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type IndexCreateResult struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index"`
}
