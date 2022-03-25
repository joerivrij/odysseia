package configs

import (
	"github.com/odysseia/plato/elastic"
)

type AlexandrosConfig struct {
	Elastic elastic.Client
	Index   string
}
