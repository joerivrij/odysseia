package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"io/ioutil"
	"strings"
)

type IndexImpl struct {
	es *elasticsearch.Client
}

func NewIndexImpl(suppliedClient *elasticsearch.Client) (*IndexImpl, error) {
	return &IndexImpl{es: suppliedClient}, nil
}

func (i *IndexImpl) CreateDocument(index string, body []byte) (*CreateResult, error) {
	var elasticResult CreateResult
	bodyString := strings.NewReader(string(body))

	esRequest := esapi.IndexRequest{
		Body:       bodyString,
		Refresh:    "true",
		Index:      index,
		DocumentID: "",
	}

	res, err := esRequest.Do(context.Background(), i.es)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	jsonBody, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (i *IndexImpl) Create(index string, request map[string]interface{}) (*IndexCreateResult, error) {
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	var elasticResult IndexCreateResult
	indexRequest := esapi.IndicesCreateRequest{
		Index: index,
		Body:  &query,
	}

	res, err := indexRequest.Do(context.Background(), i.es)
	if err != nil {
		return &elasticResult, err
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Debugf("[%s]", res.Status())
	}

	jsonBody, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = UnmarshalIndexCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (i *IndexImpl) Delete(index string) (bool, error) {
	glg.Warnf("deleting index: %s", index)

	res, err := i.es.Indices.Delete([]string{index})
	if err != nil {
		return false, err
	}

	if res.IsError() {
		if res.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return false, err
	}

	return r["acknowledged"].(bool), nil
}
