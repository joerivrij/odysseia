package elastic

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io/ioutil"
	"time"
)

type QueryImpl struct {
	es *elasticsearch.Client
}

func NewQueryImpl(suppliedClient *elasticsearch.Client) (*QueryImpl, error) {
	return &QueryImpl{es: suppliedClient}, nil
}

func (q *QueryImpl) Match(index string, request map[string]interface{}) (*Response, error) {
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	res, err := q.es.Search(
		q.es.Search.WithContext(context.Background()),
		q.es.Search.WithIndex(index),
		q.es.Search.WithBody(&query),
		q.es.Search.WithTrackTotalHits(true),
		q.es.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	return q.parseResponse(res)
}

func (q *QueryImpl) MatchWithSort(index, direction, sortField string, size int, request map[string]interface{}) (*Response, error) {
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	res, err := q.es.Search(
		q.es.Search.WithContext(context.Background()),
		q.es.Search.WithIndex(index),
		q.es.Search.WithBody(&query),
		q.es.Search.WithSize(size),
		q.es.Search.WithTrackTotalHits(true),
		q.es.Search.WithSort(fmt.Sprintf("%s:%s", sortField, direction), "mode:max"),
		q.es.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	return q.parseResponse(res)
}

func (q *QueryImpl) MatchWithScroll(index string, request map[string]interface{}) (*Response, error) {
	var elasticResult Response

	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	res, err := q.es.Search(
		q.es.Search.WithContext(context.Background()),
		q.es.Search.WithIndex(index),
		q.es.Search.WithBody(&query),
		q.es.Search.WithSize(10),
		q.es.Search.WithScroll(5*time.Second),
		q.es.Search.WithTrackTotalHits(true),
		q.es.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	firstResponse, err := q.parseResponse(res)
	if err != nil {
		return nil, err
	}

	scrollID := firstResponse.ScrollId

	for _, hit := range firstResponse.Hits.Hits {
		elasticResult.Hits.Hits = append(elasticResult.Hits.Hits, hit)
	}

	if len(firstResponse.Hits.Hits) < 10 {
		return &elasticResult, nil
	}

	for {
		scrollRes, err := q.es.Scroll(q.es.Scroll.WithScrollID(scrollID), q.es.Scroll.WithScroll(5*time.Second))
		if err != nil {
			return nil, err
		}
		defer scrollRes.Body.Close()

		if scrollRes.IsError() {
			return nil, fmt.Errorf("elasticSearch returned an error: %s", scrollRes.Status())
		}

		scrollBody, _ := ioutil.ReadAll(scrollRes.Body)
		scrollResponse, err := UnmarshalResponse(scrollBody)
		if err != nil {
			return nil, err
		}

		if len(scrollResponse.Hits.Hits) == 0 {
			break
		}

		for _, hit := range scrollResponse.Hits.Hits {
			elasticResult.Hits.Hits = append(elasticResult.Hits.Hits, hit)
		}
	}
	return &elasticResult, nil
}

func (q *QueryImpl) MatchAggregate(index string, request map[string]interface{}) (*Aggregations, error) {
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	res, err := q.es.Search(
		q.es.Search.WithContext(context.Background()),
		q.es.Search.WithIndex(index),
		q.es.Search.WithBody(&query),
		q.es.Search.WithTrackTotalHits(true),
		q.es.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	return q.parseAggregate(res)
}

func (q *QueryImpl) parseResponse(res *esapi.Response) (*Response, error) {
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	result, err := UnmarshalResponse(body)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (q *QueryImpl) parseAggregate(res *esapi.Response) (*Aggregations, error) {
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	result, err := UnmarshalAggregations(body)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
