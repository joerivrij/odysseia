package configs

import (
	"github.com/odysseia/plato/elastic"
)

type DemokritosConfig struct {
	Index      string
	SearchWord string
	Created    int
	Elastic    elastic.Client
}
