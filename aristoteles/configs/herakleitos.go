package configs

import (
	"github.com/odysseia-greek/plato/elastic"
)

type HerakleitosConfig struct {
	Index   string
	Created int
	Elastic elastic.Client
}
