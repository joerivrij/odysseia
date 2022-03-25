package configs

import (
	"github.com/odysseia/plato/elastic"
)

type HerodotosConfig struct {
	Elastic elastic.Client
	Index   string
}
