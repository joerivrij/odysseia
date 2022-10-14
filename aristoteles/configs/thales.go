package configs

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/queue"
)

type ThalesConfig struct {
	Index   string
	Created int
	Channel string
	Queue   queue.Client
	Elastic elastic.Client
}
