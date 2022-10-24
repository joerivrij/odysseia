package configs

import (
	"github.com/odysseia-greek/plato/elastic"
)

type AnaximanderConfig struct {
	Index   string
	Created int
	Elastic elastic.Client
}
