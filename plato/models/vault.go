package models

import "encoding/json"

func UnmarshalClusterKeys(data []byte) (ClusterKeys, error) {
	var r ClusterKeys
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ClusterKeys) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ClusterKeys struct {
	UnsealKeysB64         []string      `json:"unseal_keys_b64"`
	UnsealKeysHex         []string      `json:"unseal_keys_hex"`
	UnsealShares          int64         `json:"unseal_shares"`
	UnsealThreshold       int64         `json:"unseal_threshold"`
	RecoveryKeysB64       []interface{} `json:"recovery_keys_b64"`
	RecoveryKeysHex       []interface{} `json:"recovery_keys_hex"`
	RecoveryKeysShares    int64         `json:"recovery_keys_shares"`
	RecoveryKeysThreshold int64         `json:"recovery_keys_threshold"`
	RootToken             string        `json:"root_token"`
}

func (r *CreateSecretRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalSecretData(data []byte) (ElasticConfigVault, error) {
	var r ElasticConfigVault
	err := json.Unmarshal(data, &r)
	return r, err
}

type CreateSecretRequest struct {
	Data ElasticConfigVault `json:"data"`
}

func (r *ElasticConfigVault) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ElasticConfigVault struct {
	Username    string `json:"elasticUsername"`
	Password    string `json:"elasticPassword"`
	ElasticCERT string `json:"elasticCert"`
}
