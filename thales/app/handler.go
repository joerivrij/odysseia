package app

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/kubemq-io/kubemq-go"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/models"
	"github.com/odysseia/plato/transform"
	"strings"
	"sync"
	"time"
)

type ThalesHandler struct {
	Config             *configs.ThalesConfig
	QueueEmptyDuration time.Duration
}

const ClientId string = "thales"
const MessageToPull int32 = 10

func (t *ThalesHandler) HandleQueue() {
	emptyQueue := make(chan bool, 1)
	go t.queueEmpty(emptyQueue, t.QueueEmptyDuration)

	go func() {
		for {
			response, err := t.Config.Queue.PullMessages(MessageToPull, ClientId)
			if err != nil {
				glg.Error(err)
			}

			err = t.handleMessages(response)
			if err != nil {
				glg.Error(err)
			}
		}
	}()

	select {

	case <-emptyQueue:
		err := t.EmptyQueue()
		glg.Error(err)
		return

	}
}

func (t *ThalesHandler) EmptyQueue() error {
	response, err := t.Config.Queue.PullRemainingMessages(ClientId)
	if err != nil {
		return err
	}

	err = t.handleMessages(response)

	return err
}

func (t *ThalesHandler) handleMessages(response *kubemq.ReceiveQueueMessagesResponse) error {
	for _, msg := range response.Messages {
		time.Sleep(50 * time.Millisecond)
		var word models.Meros
		err := json.Unmarshal(msg.Body, &word)
		if err != nil {
			return err
		}
		found, err := t.queryWord(word)
		if err != nil {
			continue
		}
		if !found {
			t.addWord(word)
		} else {
			glg.Infof("word: %s already in dictionary", string(msg.Body))
		}
		glg.Infof("MessageID: %s, Body: %s", msg.MessageID, string(msg.Body))
	}

	return nil
}

func (t *ThalesHandler) queueEmpty(queueChannel chan bool, duration time.Duration) {
	queueStatus, _ := t.Config.Queue.Info()

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
		time.Sleep(duration)
		queueStatus, _ = t.Config.Queue.Info()
	}
}

func (t *ThalesHandler) queryWord(word models.Meros) (bool, error) {
	found := false

	strippedWord := transform.RemoveAccents(word.Greek)

	term := "greek"
	query := t.Config.Elastic.Builder().MatchQuery(term, strippedWord)
	response, err := t.Config.Elastic.Query().Match(t.Config.Index, query)

	if err != nil {
		glg.Error(err)
		return found, err
	}

	if len(response.Hits.Hits) >= 1 {
		found = true
	}

	var parsedEnglishWord string
	pronouns := []string{"a", "an", "the"}
	splitEnglish := strings.Split(word.English, " ")
	numberOfWords := len(splitEnglish)
	if numberOfWords > 1 {
		for _, pronoun := range pronouns {
			if splitEnglish[0] == pronoun {
				toJoin := splitEnglish[1:numberOfWords]
				parsedEnglishWord = strings.Join(toJoin, " ")
				break
			} else {
				parsedEnglishWord = word.English
			}
		}
	} else {
		parsedEnglishWord = word.English
	}

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if meros.English == parsedEnglishWord || meros.English == word.English {
			return true, nil
		} else {
			found = false
		}
	}

	return found, nil
}

func (t *ThalesHandler) addWord(word models.Meros) {
	var innerWaitGroup sync.WaitGroup
	jsonifiedLogos, _ := word.Marshal()
	_, err := t.Config.Elastic.Index().CreateDocument(t.Config.Index, jsonifiedLogos)

	if err != nil {
		glg.Error(err)
		return
	} else {
		innerWaitGroup.Add(1)
		go t.transformWord(word, &innerWaitGroup)
		glg.Debugf("created root word: %s", word.Greek)
		t.Config.Created++
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
	_, err := t.Config.Elastic.Index().CreateDocument(t.Config.Index, jsonifiedLogos)

	if err != nil {
		glg.Error(err)
		return
	}

	t.Config.Created++

	return
}
