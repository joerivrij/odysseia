package configs

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type DionysosConfig struct {
	ElasticClient   elasticsearch.Client
	DictionaryIndex string
	Index           string
}
