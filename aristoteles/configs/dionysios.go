package configs

import (
	"github.com/odysseia-greek/plato/cache"
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/models"
)

type DionysiosConfig struct {
	Elastic          elastic.Client
	Cache            cache.Client
	Index            string
	SecondaryIndex   string
	DeclensionConfig models.DeclensionConfig
}
