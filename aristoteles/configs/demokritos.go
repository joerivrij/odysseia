package configs

import "github.com/elastic/go-elasticsearch/v7"

type DemokritosConfig struct {
	Index         string
	SearchWord    string
	Created       int
	ElasticClient elasticsearch.Client
}
