package elastic

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"os"
)

const serviceDefault = "http://127.0.0.1:9200"
const usernameDefault = "odysseia"
const passwordDefault = "odysseia"

func CreateElasticClientFromEnvVariables() (*elasticsearch.Client, error) {
	elasticService := os.Getenv("ELASTIC_SEARCH_SERVICE")
	if elasticService == "" {
		glg.Debugf("setting ELASTIC_SEARCH_SERVICE to default: %s", serviceDefault)
		elasticService = serviceDefault
	}
	elasticUser := os.Getenv("ELASTIC_SEARCH_USER")
	if elasticUser == "" {
		glg.Debugf("setting ELASTIC_SEARCH_USER to default: %s", usernameDefault)
		elasticUser = usernameDefault
	}
	elasticPassword := os.Getenv("ELASTIC_SEARCH_PASSWORD")
	if elasticPassword == "" {
		glg.Debugf("setting ELASTIC_SEARCH_PASSWORD to default: %s", passwordDefault)
		elasticPassword = passwordDefault
	}

	glg.Debugf("%s : %s", "ELASTIC_SEARCH_PASSWORD", elasticPassword)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_USER", elasticUser)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_SERVICE", elasticService)

	es, err := CreateElasticClient(elasticPassword, elasticUser, []string{elasticService})
	if err != nil {
		return nil, err
	}

	return es, nil
}
