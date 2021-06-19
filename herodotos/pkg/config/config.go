package config

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/ianschenck/envflag"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"log"
)

type HerodotosConfig struct {
	ElasticClient elasticsearch.Client
	AuthorIndex   string
}

func Get() *HerodotosConfig {
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

	healthy := elastic.CheckHealthyStatusElasticSearch(es, 200)
	if !healthy {
		glg.Fatal("death has found me")
	}

	config := &HerodotosConfig{
		ElasticClient: *es,
		AuthorIndex:   "authors",
	}

	return config
}
