package app

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"strings"
	"sync"
)

type HerakleitosHandler struct {
	Config *configs.HerakleitosConfig
}

func (h *HerakleitosHandler) DeleteIndexAtStartUp() {
	elastic.DeleteIndex(&h.Config.ElasticClient, h.Config.Index)
}

func (h *HerakleitosHandler) CreateIndexAtStartup() {
	var buf bytes.Buffer
	indexMapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(indexMapping); err != nil {
		glg.Fatalf("Error encoding query: %s", err)
	}

	indexRequest := esapi.IndicesCreateRequest{
		Index: h.Config.Index,
		Body:  &buf,
	}
	// Perform the request with the client.
	res, err := indexRequest.Do(context.Background(), &h.Config.ElasticClient)
	if err != nil {
		glg.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		glg.Debugf("[%s]", res.Status())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			glg.Errorf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			glg.Info("created index: %s", r)
		}
	}
}

func (h *HerakleitosHandler) Add(rhema models.Rhema, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, word := range rhema.Rhemai {
		jsonifiedLogos, _ := word.Marshal()
		esRequest := esapi.IndexRequest{
			Body:       strings.NewReader(string(jsonifiedLogos)),
			Refresh:    "true",
			Index:      h.Config.Index,
			DocumentID: "",
		}

		// Perform the request with the client.
		res, err := esRequest.Do(context.Background(), &h.Config.ElasticClient)
		if err != nil {
			glg.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			glg.Debugf("[%s]", res.Status())
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				glg.Errorf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and indexed document version.

				glg.Debugf("created root word: %s", word.Greek)
				h.Config.Created++
			}
		}
	}
}
