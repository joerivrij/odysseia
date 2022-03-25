package configs

import (
	"github.com/odysseia/plato/elastic"
)

type HerakleitosConfig struct {
	Index   string
	Created int
	Elastic elastic.Client
}
