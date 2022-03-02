package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"github.com/kubemq-io/kubemq-go"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"strings"
	"sync"
)

type ParmenidessHandler struct {
	Config *configs.ParmenidesConfig
	Mq     *kubemq.QueuesClient
}

func (p *ParmenidessHandler) DeleteIndexAtStartUp() {
	elastic.DeleteIndex(&p.Config.ElasticClient, p.Config.Index)
}

func (p *ParmenidessHandler) CreateIndexAtStartup() {
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
		Index: p.Config.Index,
		Body:  &buf,
	}
	// Perform the request with the client.
	res, err := indexRequest.Do(context.Background(), &p.Config.ElasticClient)
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

func (p *ParmenidessHandler) Add(logoi models.Logos, wg *sync.WaitGroup, method, category string, queue bool) {
	defer wg.Done()
	for _, word := range logoi.Logos {
		if queue {
			p.Queue(word)
		}
		word.Category = category
		word.Method = method
		jsonifiedLogos, _ := word.Marshal()
		esRequest := esapi.IndexRequest{
			Body:       strings.NewReader(string(jsonifiedLogos)),
			Refresh:    "true",
			Index:      p.Config.Index,
			DocumentID: "",
		}

		// Perform the request with the client.
		res, err := esRequest.Do(context.Background(), &p.Config.ElasticClient)
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
				p.Config.Created++
			}
		}
	}
}

func (p *ParmenidessHandler) Queue(word models.Word) {
	meros := models.Meros{
		Greek:      word.Greek,
		English:    word.Translation,
		LinkedWord: "",
		Original:   word.Greek,
	}

	marshalled, _ := meros.Marshal()
	message := kubemq.NewQueueMessage().
		SetChannel(p.Config.Channel).SetBody(marshalled)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := p.Mq.Send(ctx, message)
	if err != nil {
		glg.Error(err.Error())
		return
	}

	glg.Infof(fmt.Sprintf("message sent to queue %s", p.Config.Channel))
}
