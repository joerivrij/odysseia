package configs

import (
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
)

type DionysosConfig struct {
	Elastic          elastic.Client
	Index            string
	SecondaryIndex   string
	DeclensionConfig models.DeclensionConfig
}
