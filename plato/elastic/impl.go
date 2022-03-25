package elastic

import (
	"crypto/x509"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"net/http"
	"time"
)

type Client interface {
	Query() Query
	Index() Index
	Builder() Builder
	Health() Health
	Access() Access
}

type Query interface {
	Match(index string, request map[string]interface{}) (*Response, error)
	MatchWithSort(index, mode, sort string, size int, request map[string]interface{}) (*Response, error)
	MatchWithScroll(index string, request map[string]interface{}) (*Response, error)
	MatchAggregate(index string, request map[string]interface{}) (*Aggregations, error)
}

type Index interface {
	CreateDocument(index string, body []byte) (*CreateResult, error)
	Create(index string, request map[string]interface{}) (*IndexCreateResult, error)
	Delete(index string) (bool, error)
}

type Builder interface {
	MatchQuery(term, queryWord string) map[string]interface{}
	MatchAll() map[string]interface{}
	MultipleMatch(mappedFields []map[string]string) map[string]interface{}
	MultiMatchWithGram(queryWord string) map[string]interface{}
	Aggregate(aggregate, field string) map[string]interface{}
	FilteredAggregate(term, queryWord, aggregate, field string) map[string]interface{}
	SearchAsYouTypeIndex(searchWord string) map[string]interface{}
	Index() map[string]interface{}
}

type Health interface {
	Check(ticks, tick time.Duration) bool
	Info() (elasticHealth models.DatabaseHealth)
}

type Access interface {
	CreateRole(name string, roleRequest CreateRoleRequest) (bool, error)
	CreateUser(name string, userCreation CreateUserRequest) (bool, error)
}

type Elastic struct {
	query   *QueryImpl
	index   *IndexImpl
	builder *BuilderImpl
	health  *HealthImpl
	access  *AccessImpl
}

func NewClient(config Config) (Client, error) {
	var err error
	var esClient *elasticsearch.Client
	if config.ElasticCERT != "" {
		esClient, err = createWithTLS(config)
		if err != nil {
			return nil, err
		}
	} else {
		esClient, err = create(config)
		if err != nil {
			return nil, err
		}
	}

	query, err := NewQueryImpl(esClient)
	if err != nil {
		return nil, err
	}

	index, err := NewIndexImpl(esClient)
	if err != nil {
		return nil, err
	}

	health, err := NewHealthImpl(esClient)
	if err != nil {
		return nil, err
	}

	access, err := NewAccessImpl(esClient)
	if err != nil {
		return nil, err
	}

	builder := NewBuilderImpl()

	es := &Elastic{query: query, index: index, builder: builder, health: health, access: access}

	return es, nil
}

func NewMockClient(fixtureFile string, statusCode int) (Client, error) {
	esClient, err := CreateMockClient(fixtureFile, statusCode)
	if err != nil {
		return nil, err
	}

	query, err := NewQueryImpl(esClient)
	if err != nil {
		return nil, err
	}

	index, err := NewIndexImpl(esClient)
	if err != nil {
		return nil, err
	}

	health, err := NewHealthImpl(esClient)
	if err != nil {
		return nil, err
	}

	access, err := NewAccessImpl(esClient)
	if err != nil {
		return nil, err
	}

	builder := NewBuilderImpl()

	es := &Elastic{query: query, index: index, builder: builder, health: health, access: access}

	return es, nil
}

func create(config Config) (*elasticsearch.Client, error) {
	glg.Info("creating elasticClient")

	cfg := elasticsearch.Config{
		Username:  config.Username,
		Password:  config.Password,
		Addresses: []string{config.Service},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		glg.Errorf("Error creating the client: %s", err)
		return nil, err
	}

	return es, nil
}

func createWithTLS(config Config) (*elasticsearch.Client, error) {
	glg.Info("creating elasticClient with tls")

	caCert := []byte(config.ElasticCERT)

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

	cfg := elasticsearch.Config{
		Username:  config.Username,
		Password:  config.Password,
		Addresses: []string{config.Service},
		Transport: tp,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		glg.Errorf("Error creating the client: %s", err)
		return nil, err
	}

	return es, nil
}

func (e *Elastic) Query() Query {
	if e == nil {
		return nil
	}
	return e.query
}

func (e *Elastic) Index() Index {
	if e == nil {
		return nil
	}
	return e.index
}

func (e *Elastic) Health() Health {
	if e == nil {
		return nil
	}
	return e.health
}

func (e *Elastic) Builder() Builder {
	if e == nil {
		return nil
	}
	return e.builder
}

func (e *Elastic) Access() Access {
	if e == nil {
		return nil
	}
	return e.access
}
