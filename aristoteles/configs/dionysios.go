package configs

import (
	"github.com/odysseia/plato/cache"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
)

type DionysiosConfig struct {
	Elastic          elastic.Client
	Cache            cache.Client
	Index            string
	SecondaryIndex   string
	DeclensionConfig models.DeclensionConfig
}
