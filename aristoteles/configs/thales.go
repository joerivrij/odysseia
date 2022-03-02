package configs

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type ThalesConfig struct {
	Index         string
	Created       int
	Channel       string
	MqAddress     string
	ElasticClient elasticsearch.Client
}
