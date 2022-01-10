package elastic

import (
	"crypto/x509"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"net/http"
)

func CreateElasticClientFromEnvVariables(config models.ElasticConfig) (*elasticsearch.Client, error) {

	glg.Debugf("%s : %s", "ELASTIC_SEARCH_PASSWORD", config.Password)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_USER", config.Username)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_SERVICE", config.Service)

	es, err := CreateElasticClient(config.Password, config.Username, []string{config.Service})
	if err != nil {
		return nil, err
	}

	return es, nil
}

func CreateElasticClientWithTlS(config models.ElasticConfig) (*elasticsearch.Client, error) {
	caCert := []byte(config.ElasticCERT)

	glg.Debugf("%s : %s", "ELASTIC_SEARCH_PASSWORD", config.Password)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_USER", config.Username)
	glg.Debugf("%s : %s", "ELASTIC_SEARCH_SERVICE", config.Service)

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

	es, err := CreateElasticClientWithTLS(config.Password, config.Username, []string{config.Service}, tp)
	if err != nil {
		return nil, err
	}

	return es, nil
}
