package configs

import (
	"github.com/odysseia-greek/plato/elastic"
)

type SokratesConfig struct {
	Elastic    elastic.Client
	SearchWord string
	Index      string
}
