package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CreateElasticClient Create an elasticclient and return a pointer to that client
func CreateElasticClient(password, username string, elasticService []string) (*elasticsearch.Client, error) {
	glg.Info("creating elasticClient")

	cfg := elasticsearch.Config{
		Username:  username,
		Password:  password,
		Addresses: elasticService,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		glg.Errorf("Error creating the client: %s", err)
		return nil, err
	}

	return es, nil
}

func CreateElasticClientWithTLS(password, username string, elasticService []string, tp http.RoundTripper) (*elasticsearch.Client, error) {
	glg.Info("creating elasticClient")

	cfg := elasticsearch.Config{
		Username:  username,
		Password:  password,
		Addresses: elasticService,
		Transport: tp,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		glg.Errorf("Error creating the client: %s", err)
		return nil, err
	}

	return es, nil
}

// CheckHealthyStatusElasticSearch Check if elastic is ready to receive requests
func CheckHealthyStatusElasticSearch(es *elasticsearch.Client, ticks time.Duration) bool {
	healthy := false

	ticker := time.NewTicker(1 * time.Second)
	timeout := time.After(ticks)

	for {
		select {
		case t := <-ticker.C:
			glg.Infof("tick: %s", t)
			res, err := es.Info()
			if err != nil {
				glg.Errorf("Error getting response: %s", err)
				continue
			}

			r, err := parseBody(res)
			if err != nil {
				continue
			}

			glg.Infof("serverVersion: %s", r["version"].(map[string]interface{})["number"])
			glg.Infof("serverName: %s", r["name"])
			glg.Infof("clusterName: %s", r["cluster_name"])
			healthy = true
			ticker.Stop()

		case <-timeout:
			ticker.Stop()
		}
		break
	}

	return healthy
}

func parseBody(res *esapi.Response) (map[string]interface{}, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			glg.Errorf("error closing elastic response body: %w", err)
		}
	}(res.Body)
	// Check response status
	if res.IsError() {
		glg.Errorf("Error: %s", res.String())
		return nil, fmt.Errorf(res.String())
	}

	var r map[string]interface{}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		glg.Errorf("Error parsing the response body: %s", err)
		return nil, err
	}

	return r, nil
}

// CheckHealth Check if elastic connection is healthy
func CheckHealth(es *elasticsearch.Client) (elasticHealth models.DatabaseHealth) {
	res, err := es.Info()

	if err != nil {
		glg.Errorf("Error getting response: %s", err)
		elasticHealth.Healthy = false
		return elasticHealth
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		elasticHealth.Healthy = false
		glg.Errorf("Error: %s", res.String())
		return elasticHealth
	}

	var r map[string]interface{}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		elasticHealth.Healthy = false
		return elasticHealth
	}

	elasticHealth.ClusterName = fmt.Sprintf("%s", r["cluster_name"])
	elasticHealth.ServerName = fmt.Sprintf("%s", r["name"])
	elasticHealth.ServerVersion = fmt.Sprintf("%s", r["version"].(map[string]interface{})["number"])
	elasticHealth.Healthy = true

	return elasticHealth
}

// CreateRole to create a role in ES
func CreateRole(elasticClient *elasticsearch.Client, name string, roleRequest models.CreateRoleRequest) (bool, error) {
	jsonRole, err := roleRequest.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonRole)
	res, _ := elasticClient.Security.PutRole(name, buffer)
	glg.Debug(res)
	return true, nil
}

// CreateUser Creates a new user
func CreateUser(elasticClient *elasticsearch.Client, name string, userCreation models.CreateUserRequest) (bool, error) {
	jsonUser, err := userCreation.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonUser)
	res, _ := elasticClient.Security.PutUser(name, buffer)
	glg.Debug(res)
	return true, nil
}

// DeleteIndex delete an index without checking for success
func DeleteIndex(es *elasticsearch.Client, index string) {
	glg.Warnf("deleting index: %s", index)

	res, err := es.Indices.Delete([]string{index})
	if err != nil {
		glg.Errorf("Error getting response: %s", err)
		return
	}

	glg.Infof("status: %s", strconv.Itoa(res.StatusCode))

	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseBody := string(responseBytes)

	switch res.StatusCode {
	case 200:
		glg.Infof("delete success: %s", responseBody)
	case 404:
		glg.Warnf("could not find index: %s", responseBody)
	default:
		glg.Errorf("something else went wrong: %s", responseBody)
	}

	return
}

func QueryWithMatch(elasticClient elasticsearch.Client, index, term, word string) (*models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				term: word,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		glg.Debug("Error encoding query: %s", err)
		return nil, err
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		glg.Debug("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Errorf("Error getting response: %s", res.Status())
		return nil, fmt.Errorf("elasticSearch returned an error: %s", res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func QueryWithMatchAll(elasticClient elasticsearch.Client, index string) (*models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		glg.Debug("Error encoding query: %s", err)
		return nil, err
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		glg.Debug("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Errorf("Error getting response: %s", res.Status())
		return nil, fmt.Errorf("elasticSearch returned an error: %s", res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func QueryMultiMatchWithGrams(elasticClient elasticsearch.Client, index, queryWord string) (*models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"size": 15,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": queryWord,
				"type":  "bool_prefix",
				"fields": [3]string{
					"greek", "greek._2gram", "greek._3gram",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		glg.Errorf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		glg.Errorf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Errorf("Error getting response: %s", res.Status())
		return nil, fmt.Errorf("elasticSearch returned an error: %s", res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func QueryOnId(elasticClient elasticsearch.Client, index, id string) (*models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": id,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		glg.Errorf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		glg.Errorf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Errorf("Error getting response: %s", res.Status())
		return nil, fmt.Errorf("elasticSearch returned an error: %s", res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func QueryWithScroll(elasticClient elasticsearch.Client, index, term, word string) (*models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				term: word,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		glg.Debug("Error encoding query: %s", err)
		return nil, err
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithSize(10),
		elasticClient.Search.WithScroll(5*time.Second),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		glg.Errorf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Errorf("Error getting response: %s", res.Status())
		return nil, fmt.Errorf("elasticSearch returned an error: %s", res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	firstResponse, err := models.UnmarshalElasticResponse(body)
	if err != nil {
		glg.Errorf("Error decoding json: %s", err)
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
		scrollRes, err := elasticClient.Scroll(elasticClient.Scroll.WithScrollID(scrollID), elasticClient.Scroll.WithScroll(5*time.Second))
		if err != nil {
			glg.Errorf("Error getting response: %s", err)
			return nil, err
		}
		defer scrollRes.Body.Close()

		if scrollRes.IsError() {
			glg.Errorf("Error getting response: %s", scrollRes.Status())
			return nil, fmt.Errorf("elasticSearch returned an error: %s", scrollRes.Status())
		}

		scrollBody, _ := ioutil.ReadAll(scrollRes.Body)
		scrollResponse, err := models.UnmarshalElasticResponse(scrollBody)
		if err != nil {
			glg.Errorf("Error decoding scrollResponse: %s", err)
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

func QueryWithDescendingSort(elasticClient elasticsearch.Client, index, sort string, results int) (*models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		glg.Errorf("Error getting response: %s", err)
		return nil, err
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithSize(results),
		elasticClient.Search.WithSort(fmt.Sprintf("%s:desc", sort), "mode:max"),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		glg.Errorf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Errorf("Error getting response: %s", res.Status())
		return nil, fmt.Errorf("elasticSearch returned an error: %s", res.Status())
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}
