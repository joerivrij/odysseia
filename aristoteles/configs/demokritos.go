package configs

import (
	"github.com/odysseia-greek/plato/elastic"
)

type DemokritosConfig struct {
	Index      string
	SearchWord string
	Created    int
	Elastic    elastic.Client
}
