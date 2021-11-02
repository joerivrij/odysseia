package app

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"time"
)

type AlexandrosConfig struct {
	ElasticClient elasticsearch.Client
	Index         string
}

func Get(ticks time.Duration, es *elasticsearch.Client) (bool, *AlexandrosConfig) {
	healthy := elastic.CheckHealthyStatusElasticSearch(es, ticks)
	if !healthy {
		glg.Errorf("elasticClient unhealthy after %s ticks", ticks)
		return healthy, nil
	}

	config := &AlexandrosConfig{
		ElasticClient: *es,
		Index:         "dictionary",
	}

	return healthy, config
}
