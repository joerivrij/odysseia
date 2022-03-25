package configs

import (
	"github.com/odysseia/plato/elastic"
)

type SokratesConfig struct {
	Elastic    elastic.Client
	SearchWord string
	Index      string
}
