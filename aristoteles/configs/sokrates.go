package configs

import "github.com/elastic/go-elasticsearch/v7"

type SokratesConfig struct {
	ElasticClient elasticsearch.Client
	SearchWord    string
}
