package elastic

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// create an elasticclient and return a pointer to that client
func CreateElasticClient(password, username string, elasticService []string) (*elasticsearch.Client, error) {
	glg.Info("creating elasticClient")
	glg.Info("creating elasticClient")

	cfg := elasticsearch.Config{
		Username: username,
		Password: password,
		Addresses: elasticService,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		glg.Fatalf("Error creating the client: %s", err)
		return nil, err
	}

	// Print client and server version numbers.
	glg.Infof("elasticClient version: %s", elasticsearch.Version)
	glg.Info(strings.Repeat("~", 37))

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