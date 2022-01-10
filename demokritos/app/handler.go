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
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"sync"
	"unicode"
)

type DemokritosHandler struct {
	Config *configs.DemokritosConfig
}

func (d *DemokritosHandler) AddDirectoryToElastic(biblos models.Biblos, wg *sync.WaitGroup) {
	defer wg.Done()
	var innerWaitGroup sync.WaitGroup
	for _, word := range biblos.Biblos {
		jsonifiedLogos, _ := word.Marshal()
		esRequest := esapi.IndexRequest{
			Body:       strings.NewReader(string(jsonifiedLogos)),
			Refresh:    "true",
			Index:      d.Config.Index,
			DocumentID: "",
		}

		// Perform the request with the client.
		res, err := esRequest.Do(context.Background(), &d.Config.ElasticClient)
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
				innerWaitGroup.Add(1)
				go d.transformWord(word, &innerWaitGroup)
				glg.Debugf("created root word: %s", word.Greek)
				d.Config.Created++
			}
		}
	}
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func (d *DemokritosHandler) DeleteIndexAtStartUp() {
	elastic.DeleteIndex(&d.Config.ElasticClient, d.Config.Index)
}

func (d *DemokritosHandler) CreateIndexAtStartup() {
	var buf bytes.Buffer
	indexMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				d.Config.SearchWord: map[string]interface{}{
					"type": "search_as_you_type",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(indexMapping); err != nil {
		glg.Fatalf("Error encoding query: %s", err)
	}

	indexRequest := esapi.IndicesCreateRequest{
		Index: d.Config.Index,
		Body:  &buf,
	}
	// Perform the request with the client.
	res, err := indexRequest.Do(context.Background(), &d.Config.ElasticClient)
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

func (d *DemokritosHandler) transformWord(m models.Meros, wg *sync.WaitGroup) {
	defer wg.Done()
	strippedWord := removeAccents(m.Greek)
	word := models.Meros{
		Greek:      strippedWord,
		English:    m.English,
		LinkedWord: m.LinkedWord,
		Original:   m.Greek,
	}

	jsonifiedLogos, _ := word.Marshal()
	esRequest := esapi.IndexRequest{
		Body:       strings.NewReader(string(jsonifiedLogos)),
		Refresh:    "true",
		Index:      d.Config.Index,
		DocumentID: "",
	}

	// Perform the request with the client.
	res, err := esRequest.Do(context.Background(), &d.Config.ElasticClient)
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
			glg.Debugf("created parsed word: %s", strippedWord)
			d.Config.Created++
		}
	}

	return
}
