package configs

import (
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/queue"
)

type ThalesConfig struct {
	Index   string
	Created int
	Channel string
	Queue   queue.Client
	Elastic elastic.Client
}
