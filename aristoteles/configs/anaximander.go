package configs

import (
	"github.com/odysseia/plato/elastic"
)

type AnaximanderConfig struct {
	Index   string
	Created int
	Elastic elastic.Client
}
