package configs

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/queue"
)

type ParmenidesConfig struct {
	Index   string
	Created int
	Elastic elastic.Client
	Queue   queue.Client
}
