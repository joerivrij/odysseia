package configs

import (
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/queue"
)

type ParmenidesConfig struct {
	Index   string
	Created int
	Elastic elastic.Client
	Queue   queue.Client
}
