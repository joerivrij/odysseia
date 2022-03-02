package configs

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type ParmenidesConfig struct {
	Index         string
	Channel       string
	MqAddress     string
	Created       int
	ElasticClient elasticsearch.Client
}
