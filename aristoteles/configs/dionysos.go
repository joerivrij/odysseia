package configs

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/odysseia/plato/models"
)

type DionysosConfig struct {
	ElasticClient    elasticsearch.Client
	Index            string
	SecondaryIndex   string
	DeclensionConfig models.DeclensionConfig
}
