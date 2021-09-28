package app

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"os"
	"time"
)

type DionysosConfig struct {
	ElasticClient elasticsearch.Client
	DeclensionHandler DeclensionHandler
}

const alexandrosDefault = "http://minikube-odysseia.test"
const versionDefault = "v1"
const apiNameDefault = "alexandros"
const endpointDefault = "search"

func Get(ticks time.Duration, es *elasticsearch.Client) (bool, *DionysosConfig) {
	healthy := elastic.CheckHealthyStatusElasticSearch(es, ticks)
	if !healthy {
		glg.Errorf("elasticClient unhealthy after %s ticks", ticks)
		return healthy, nil
	}

	alexandrosBaseUrl := os.Getenv("ALEXANDROS_BASE_URL")
	if alexandrosBaseUrl == "" {
		glg.Debugf("setting ALEXANDROS_BASE_URL to default: %s", alexandrosDefault)
		alexandrosBaseUrl = alexandrosDefault
	}

	version := os.Getenv("VERSION")
	if alexandrosBaseUrl == "" {
		glg.Debugf("setting VERSION to default: %s", versionDefault)
		version = versionDefault
	}

	apiName := os.Getenv("API_NAME")
	if alexandrosBaseUrl == "" {
		glg.Debugf("setting API_NAME to default: %s", apiNameDefault)
		apiName = apiNameDefault
	}

	searchEndpoint := os.Getenv("SEARCH_ENDPOINT")
	if alexandrosBaseUrl == "" {
		glg.Debugf("setting SEARCH_ENDPOINT to default: %s", endpointDefault)
		searchEndpoint = endpointDefault
	}

	config := &DionysosConfig{
		ElasticClient: *es,
		DeclensionHandler: DeclensionHandler{
			BaseUrl: alexandrosBaseUrl,
			Version: version,
			ApiName: apiName,
			SearchWordEndPoint: searchEndpoint,
			Index: "dionysos",
			ElasticClient: *es,
		},
	}

	return healthy, config
}
