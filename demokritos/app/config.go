package app

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
)

type DemokritosConfig struct {
	Index         string
	SearchWord    string
	Created       int
	ElasticClient elasticsearch.Client
}

func Get() *DemokritosConfig {
	index := "dictionary"
	searchWord := "greek"

	confManager, err := aristoteles.NewConfig()
	if err != nil {
		glg.Error(err)
		glg.Fatal("unable to fetch configuration")
	}

	elasticClient, err := confManager.GetElasticClient()
	if err != nil {
		glg.Fatal("failed to create client")
	}

	return &DemokritosConfig{
		Index:         index,
		SearchWord:    searchWord,
		Created:       0,
		ElasticClient: *elasticClient,
	}
}
