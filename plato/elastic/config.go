package elastic

import (
	"crypto/x509"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"net/http"
	"os"
)

const serviceDefault = "http://localhost:9200"
const serviceDefaultTlS = "https://localhost:9200"
const usernameDefault = "elastic"
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

func CreateElasticClientWithTlS(config models.ElasticConfigVault) (*elasticsearch.Client, error) {
	elasticService := os.Getenv("ELASTIC_SEARCH_SERVICE")
	if elasticService == "" {
		glg.Debugf("setting ELASTIC_SEARCH_SERVICE to default: %s", serviceDefaultTlS)
		elasticService = serviceDefaultTlS
	}

	caCert := []byte(config.ElasticCERT)

	glg.Debugf("%s : %s", "ELASTIC_SEARCH_PASSWORD", config.Password)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_USER", config.Username)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_SERVICE", elasticService)

	// --> Clone the default HTTP transport

	tp := http.DefaultTransport.(*http.Transport).Clone()

	// --> Initialize the set of root certificate authorities
	//
	var err error

	if tp.TLSClientConfig.RootCAs, err = x509.SystemCertPool(); err != nil {
		glg.Fatalf("ERROR: Problem adding system CA: %s", err)
	}

	// --> Add the custom certificate authority
	//
	if ok := tp.TLSClientConfig.RootCAs.AppendCertsFromPEM(caCert); !ok {
		glg.Fatalf("ERROR: Problem adding CA from file %q", caCert)
	}

	es, err := CreateElasticClientWithTLS(config.Password, config.Username, []string{elasticService}, tp)
	if err != nil {
		return nil, err
	}

	return es, nil
}


func CreateElasticClientFromEnvVariablesWithTLS(caCert []byte) (*elasticsearch.Client, error) {
	elasticService := os.Getenv("ELASTIC_SEARCH_SERVICE")
	if elasticService == "" {
		glg.Debugf("setting ELASTIC_SEARCH_SERVICE to default: %s", serviceDefaultTlS)
		elasticService = serviceDefaultTlS
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

	// --> Clone the default HTTP transport

	tp := http.DefaultTransport.(*http.Transport).Clone()

	// --> Initialize the set of root certificate authorities
	//
	var err error

	if tp.TLSClientConfig.RootCAs, err = x509.SystemCertPool(); err != nil {
		glg.Fatalf("ERROR: Problem adding system CA: %s", err)
	}

	// --> Add the custom certificate authority
	//
	if ok := tp.TLSClientConfig.RootCAs.AppendCertsFromPEM(caCert); !ok {
		glg.Fatalf("ERROR: Problem adding CA from file %q", caCert)
	}

	es, err := CreateElasticClientWithTLS(elasticPassword, elasticUser, []string{elasticService}, tp)
	if err != nil {
		return nil, err
	}

	return es, nil
}
