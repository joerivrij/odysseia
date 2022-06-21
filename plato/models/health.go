package models

import (
	"encoding/json"
)

func UnmarshalHealth(data []byte) (Health, error) {
	var r Health
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Health) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Health struct {
	Healthy  bool           `json:"healthy"`
	Time     string         `json:"time"`
	Database DatabaseHealth `json:"databaseHealth"`
	Memory   Memory         `json:"memory"`
}

type DatabaseHealth struct {
	Healthy       bool   `json:"healthy"`
	ClusterName   string `json:"clusterName,omitempty"`
	ServerName    string `json:"serverName,omitempty"`
	ServerVersion string `json:"serverVersion,omitempty"`
}

type Memory struct {
	Free       uint64 `json:"free"`
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"totalAlloc"`
	Sys        uint64 `json:"sys"`
	Unit       string `json:"unit"`
}
