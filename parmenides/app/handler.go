package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia/aristoteles/configs"
	"sync"
)

type ParmenidesHandler struct {
	Config *configs.ParmenidesConfig
}

func (p *ParmenidesHandler) DeleteIndexAtStartUp() error {
	deleted, err := p.Config.Elastic.Index().Delete(p.Config.Index)
	glg.Infof("deleted index: %s %v", p.Config.Index, deleted)
	if err != nil {
		return err
	}

	return nil
}

func (p *ParmenidesHandler) CreateIndexAtStartup() error {
	indexMapping := p.Config.Elastic.Builder().Index()
	created, err := p.Config.Elastic.Index().Create(p.Config.Index, indexMapping)
	if err != nil {
		return err
	}

	glg.Infof("created index: %s %v", created.Index, created.Acknowledged)

	return nil
}

func (p *ParmenidesHandler) Add(logoi models.Logos, wg *sync.WaitGroup, method, category string, queue bool) error {
	defer wg.Done()
	for _, word := range logoi.Logos {
		if queue {
			meros := models.Meros{
				Greek:      word.Greek,
				English:    word.Translation,
				LinkedWord: "",
				Original:   word.Greek,
			}

			marshalled, _ := meros.Marshal()

			err := p.Queue(marshalled)
			if err != nil {
				glg.Error(err)
			}
		}

		if method == "mouseion" {
			glg.Debug(method)
		}
		word.Category = category
		word.Method = method
		jsonifiedLogos, _ := word.Marshal()
		_, err := p.Config.Elastic.Index().CreateDocument(p.Config.Index, jsonifiedLogos)

		if err != nil {
			return err
		}

		glg.Infof("created word: %s with translation %s | method: %s | category: %s", word.Greek, word.Translation, word.Method, word.Category)

		p.Config.Created++
	}
	return nil
}

func (p *ParmenidesHandler) Queue(marshalled []byte) error {
	res, err := p.Config.Queue.SendQueueMessage(marshalled)
	if err != nil {
		return err
	}

	message := string(marshalled)
	glg.Debugf("message queued: %s with payload: %s", res.MessageID, message)
	return nil
}
