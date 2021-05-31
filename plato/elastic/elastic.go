package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/lexiko/plato/models"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

// create an elasticclient and return a pointer to that client
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

// Check if elastic is ready to receive requests
func CheckHealthyStatusElasticSearch(es *elasticsearch.Client, ticks time.Duration) bool {
	healthy := false

	ticker := time.NewTicker(1 * time.Second)
	timeout := time.After(ticks * time.Second)

	for {
		select {
		case t := <-ticker.C:
			glg.Infof("tick: %s", t)
			res, err := es.Info()
			if err != nil {
				glg.Errorf("Error getting response: %s", err)
				continue
			}
			defer res.Body.Close()
			// Check response status
			if res.IsError() {
				glg.Errorf("Error: %s", res.String())
			}

			var r map[string]interface{}

			// Deserialize the response into a map.
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				glg.Errorf("Error parsing the response body: %s", err)
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

// delete an index without checking for success
func DeleteIndex(es *elasticsearch.Client, index string) {
	glg.Warnf("deleting index: %s", index)

	res, err := es.Indices.Delete([]string{index})
	if err != nil {
		glg.Errorf("Error getting response: %s", err)
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

func QueryWithMatchAll(elasticClient elasticsearch.Client, index string) (models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			glg.Error(err)
		} else {
			// Print the response status and error information.
			glg.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return elasticResult, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return elasticResult, err
	}

	return elasticResult, nil
}

func QueryMultiMatchWithGrams(elasticClient elasticsearch.Client, index, queryWord string) (models.ElasticResponse, error) {
	var elasticResult models.ElasticResponse
	var buf bytes.Buffer
	query := map[string]interface{}{
		"size": 10,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": queryWord,
				"type": "bool_prefix",
				"fields": [3]string{
					"greek", "greek._2gram", "greek._3gram",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			glg.Error(err)
		} else {
			// Print the response status and error information.
			glg.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return elasticResult, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return elasticResult, err
	}

	return elasticResult, nil
}

func QueryOnId(elasticClient elasticsearch.Client, index, id string) (models.ElasticResponse, error) {
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
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			glg.Error(err)
		} else {
			// Print the response status and error information.
			glg.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return elasticResult, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalElasticResponse(body)
	if err != nil {
		return elasticResult, err
	}

	return elasticResult, nil
}
