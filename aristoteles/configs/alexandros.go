package configs

import "github.com/elastic/go-elasticsearch/v7"

type AlexandrosConfig struct {
	ElasticClient elasticsearch.Client
	Index         string
}
