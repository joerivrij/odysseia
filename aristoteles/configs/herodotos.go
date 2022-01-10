package configs

import "github.com/elastic/go-elasticsearch/v7"

type HerodotosConfig struct {
	ElasticClient elasticsearch.Client
	Index         string
}
