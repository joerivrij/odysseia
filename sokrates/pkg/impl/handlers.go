package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"net/http"
	"sokrates/pkg/config"
	"sokrates/pkg/middleware"
	"sokrates/pkg/models"
	"strings"
)

type SokratesHandler struct {
	Config *config.SokratesConfig
}

// PingPong pongs the ping
func (s *SokratesHandler)PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (s *SokratesHandler)CreateDocuments(w http.ResponseWriter, req *http.Request) {
	var logosModel = models.Logos{Logos:
		[]models.Word{
			{
				Greek:   "ἀγαθός",
				Dutch:   "goed",
				Chapter: 0,
			},
			{
				Greek:   "τὸ αγαθών",
				Dutch:   "het goede",
				Chapter: 0,
		},
	}}

	for _, word := range logosModel.Logos {
		jsonifiedLogos, _ := word.Marshal()
		esRequest := esapi.IndexRequest{
			Body:        strings.NewReader(string(jsonifiedLogos)),
			Refresh:    "true",
			Index:      "logos",
			DocumentID: middleware.CreateGUID(),
		}

		// Perform the request with the client.
		res, err := esRequest.Do(context.Background(), &s.Config.ElasticClient)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s]", res.Status())
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and indexed document version.
				log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}
	}



	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (s *SokratesHandler)QueryAllForIndex(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"greek": "ἀγαθός",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := s.Config.ElasticClient.Search(
		s.Config.ElasticClient.Search.WithContext(context.Background()),
		s.Config.ElasticClient.Search.WithIndex("logos"),
		s.Config.ElasticClient.Search.WithBody(&buf),
		s.Config.ElasticClient.Search.WithTrackTotalHits(true),
		s.Config.ElasticClient.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r  map[string]interface{}

	var someArray []interface{}
	json.NewDecoder(res.Body).Decode(&r)
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		tie := hit.(map[string]interface{})["_source"]
		someArray = append(someArray, tie)
	}

	middleware.ResponseWithJson(w, someArray)


}