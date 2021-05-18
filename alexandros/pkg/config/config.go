package config

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/ianschenck/envflag"
	"github.com/kpango/glg"
	"github.com/lexiko/plato/elastic"
	"log"
	"time"
)

type AlexandrosConfig struct {
	ElasticClient elasticsearch.Client
}

func Get(ticks time.Duration) (bool, *AlexandrosConfig) {
	elasticService := envflag.String("ELASTIC_SEARCH_SERVICE", "http://127.0.0.1:9200", "location of the es service")
	elasticUser := envflag.String("ELASTIC_SEARCH_USER", "sokrates", "es username")
	elasticPassword := envflag.String("ELASTIC_SEARCH_PASSWORD", "sokrates", "es password")

	envflag.Parse()

	glg.Debugf("%s : %s", "ELASTIC_SEARCH_PASSWORD", *elasticPassword)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_USER", *elasticUser)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_SERVICE", *elasticService)

	es, err := elastic.CreateElasticClient(*elasticPassword, *elasticUser, []string{*elasticService})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	healthy := elastic.CheckHealthyStatusElasticSearch(es, ticks)
	if !healthy {
		glg.Errorf("elasticClient unhealthy after %s ticks", ticks)
		return healthy, nil
	}

	config := &AlexandrosConfig{
		ElasticClient: *es,
	}

	return healthy, config
}
