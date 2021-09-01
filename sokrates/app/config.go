package app

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"time"
)

type SokratesConfig struct {
	ElasticClient elasticsearch.Client
	SearchTerm    string
}

func Get(ticks time.Duration, es *elasticsearch.Client) (bool, *SokratesConfig) {

	healthy := elastic.CheckHealthyStatusElasticSearch(es, ticks)
	if !healthy {
		glg.Errorf("elasticClient unhealthy after %s ticks", ticks)
		return healthy, nil
	}

	config := &SokratesConfig{
		ElasticClient: *es,
		SearchTerm:    "greek",
	}

	return healthy, config
}
