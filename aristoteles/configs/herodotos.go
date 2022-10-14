package configs

import (
	"github.com/odysseia-greek/plato/elastic"
)

type HerodotosConfig struct {
	Elastic elastic.Client
	Index   string
}
