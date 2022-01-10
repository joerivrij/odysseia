package configs

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/odysseia/plato/models"
)

type DionysosConfig struct {
	ElasticClient    elasticsearch.Client
	DictionaryIndex  string
	Index            string
	DeclensionConfig models.DeclensionConfig
}
