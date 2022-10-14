package configs

import (
	"github.com/odysseia-greek/plato/elastic"
)

type AlexandrosConfig struct {
	Elastic elastic.Client
	Index   string
}
