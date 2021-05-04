package config

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/ianschenck/envflag"
	"github.com/kpango/glg"
	"log"
	"strings"
)

type SokratesConfig struct {
	ElasticService string
	ElastictUser string
	ElasticPassword string
	ElasticClient elasticsearch.Client
	SearchTerm string
}

func Get() *SokratesConfig {
	elasticService := envflag.String("ELASTIC_SEARCH_SERVICE","http://127.0.0.1:9200", "location of the es service")
	elasticUser := envflag.String("ELASTIC_SEARCH_USER","sokrates", "es username")
	elasticPassword := envflag.String("ELASTIC_SEARCH_PASSWORD","sokrates", "es password")

	envflag.Parse()

	glg.Debugf("%s : %s", "ELASTIC_SEARCH_PASSWORD", *elasticPassword)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_USER", *elasticUser)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_SERVICE", *elasticService)


	cfg := elasticsearch.Config{
		Username: "elastic",
		Password: "changeme",
		Addresses: []string{
			*elasticService,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	var r map[string]interface{}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	config := &SokratesConfig{
		ElasticService: *elasticService,
		ElastictUser: *elasticUser,
		ElasticPassword: *elasticPassword,
		ElasticClient: *es,
		SearchTerm: "greek",
	}

	return config
}