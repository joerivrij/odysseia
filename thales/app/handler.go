package app

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"github.com/kubemq-io/kubemq-go"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/transform"
	"strings"
	"sync"
	"time"
)

type ThalesHandler struct {
	Config *configs.ThalesConfig
	Mq     *kubemq.QueuesClient
}

func (t *ThalesHandler) HandleQueue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emptyQueue := make(chan bool, 1)
	go t.queueEmpty(emptyQueue)
	done, _ := t.Mq.Subscribe(ctx, &kubemq.ReceiveQueueMessagesRequest{
		ClientID:            "thales",
		Channel:             t.Config.Channel,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     2,
	}, func(response *kubemq.ReceiveQueueMessagesResponse, err error) {
		if err != nil {
			glg.Error(err)
		}
		for _, msg := range response.Messages {
			time.Sleep(50 * time.Millisecond)
			var word models.Meros
			err := json.Unmarshal(msg.Body, &word)
			if err != nil {
				glg.Error(err)
			}
			found := t.queryWord(word)
			if !found {
				t.addWord(word)
			}
			glg.Infof("MessageID: %s, Body: %s", msg.MessageID, string(msg.Body))

		}
	})

	select {

	case <-emptyQueue:
		done <- struct{}{}
	}
	done <- struct{}{}

	return
}

func (t *ThalesHandler) queueEmpty(queueChannel chan bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queueStatus, _ := t.Mq.QueuesInfo(ctx, t.Config.Channel)

	for {
		for _, queue := range queueStatus.Queues {
			if queue.Name == t.Config.Channel {
				glg.Debugf("%v messages remaining", queue.Waiting)
				if queue.Waiting == 0 {
					queueChannel <- true
				}
			}
		}
		//check every 10 seconds for an empty queue
		time.Sleep(10000 * time.Millisecond)
		queueStatus, _ = t.Mq.QueuesInfo(ctx, t.Config.Channel)
	}
}

func (t *ThalesHandler) queryWord(word models.Meros) bool {
	found := false

	strippedWord := transform.RemoveAccents(word.Greek)

	term := "greek"
	response, err := elastic.QueryWithMatch(t.Config.ElasticClient, t.Config.Index, term, strippedWord)

	if err != nil {
		glg.Error(err)
	}

	if len(response.Hits.Hits) >= 1 {
		found = true
	}

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if meros.English != word.English {
			found = false
		}
	}

	return found
}

func (t *ThalesHandler) addWord(word models.Meros) {
	var innerWaitGroup sync.WaitGroup
	jsonifiedLogos, _ := word.Marshal()
	esRequest := esapi.IndexRequest{
		Body:       strings.NewReader(string(jsonifiedLogos)),
		Refresh:    "true",
		Index:      t.Config.Index,
		DocumentID: "",
	}

	// Perform the request with the client.
	res, err := esRequest.Do(context.Background(), &t.Config.ElasticClient)
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
			go t.transformWord(word, &innerWaitGroup)
			glg.Debugf("created root word: %s", word.Greek)
			t.Config.Created++
		}
	}
}

func (t *ThalesHandler) transformWord(word models.Meros, wg *sync.WaitGroup) {
	defer wg.Done()
	strippedWord := transform.RemoveAccents(word.Greek)
	meros := models.Meros{
		Greek:      strippedWord,
		English:    word.English,
		LinkedWord: word.LinkedWord,
		Original:   word.Greek,
	}

	jsonifiedLogos, _ := meros.Marshal()
	esRequest := esapi.IndexRequest{
		Body:       strings.NewReader(string(jsonifiedLogos)),
		Refresh:    "true",
		Index:      t.Config.Index,
		DocumentID: "",
	}

	// Perform the request with the client.
	res, err := esRequest.Do(context.Background(), &t.Config.ElasticClient)
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
			t.Config.Created++
		}
	}

	return
}
