package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"log"
	"time"
)

func QueryWithMatch(elasticClient elasticsearch.Client, index, term, word string) (models.Logos, map[string]interface{}) {
	var queryResult models.Logos
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				term: word,
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

		return queryResult, e
	}

	var r map[string]interface{}

	json.NewDecoder(res.Body).Decode(&r)
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		//glg.Debugf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		tie, _ := json.Marshal(hit.(map[string]interface{})["_source"])
		queryWord, _ := models.UnmarshalWord(tie)
		queryResult.Logos = append(queryResult.Logos, queryWord)
	}

	return queryResult, nil
}

func QueryWithScroll(elasticClient elasticsearch.Client, index, term, word string) (models.Logos, map[string]interface{}) {
	var queryResult models.Logos
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				term: word,
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
		elasticClient.Search.WithSize(10),
		elasticClient.Search.WithScroll(5*time.Second),
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

		return queryResult, e
	}

	var r map[string]interface{}

	json.NewDecoder(res.Body).Decode(&r)

	scrollID := fmt.Sprintf("%v", r["_scroll_id"])
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		//glg.Debugf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		tie, _ := json.Marshal(hit.(map[string]interface{})["_source"])
		queryWord, _ := models.UnmarshalWord(tie)
		queryResult.Logos = append(queryResult.Logos, queryWord)
	}

	for {
		scrollRes, err := elasticClient.Scroll(elasticClient.Scroll.WithScrollID(scrollID), elasticClient.Scroll.WithScroll(5*time.Second))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		defer scrollRes.Body.Close()

		if scrollRes.IsError() {
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

			return queryResult, e
		}

		var scroll map[string]interface{}

		json.NewDecoder(scrollRes.Body).Decode(&scroll)

		if len(scroll["hits"].(map[string]interface{})["hits"].([]interface{})) == 0 {
			break
		}
		for _, hit := range scroll["hits"].(map[string]interface{})["hits"].([]interface{}) {
			glg.Debugf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
			tie, _ := json.Marshal(hit.(map[string]interface{})["_source"])
			queryWord, _ := models.UnmarshalWord(tie)
			queryResult.Logos = append(queryResult.Logos, queryWord)
		}
	}
	return queryResult, nil
}

func QueryLastChapter(elasticClient elasticsearch.Client, index string) (int64, map[string]interface{}) {
	var queryResult models.Logos
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
		elasticClient.Search.WithSize(1),
		elasticClient.Search.WithSort("chapter:desc", "mode:max"),
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

		return 0, e
	}

	var r map[string]interface{}

	json.NewDecoder(res.Body).Decode(&r)
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		glg.Debugf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		tie, _ := json.Marshal(hit.(map[string]interface{})["_source"])
		queryWord, _ := models.UnmarshalWord(tie)
		queryResult.Logos = append(queryResult.Logos, queryWord)
	}

	return queryResult.Logos[0].Chapter, nil
}
