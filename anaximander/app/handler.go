package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia/aristoteles/configs"
	"sync"
)

type AnaximanderHandler struct {
	Config *configs.AnaximanderConfig
}

func (a *AnaximanderHandler) DeleteIndexAtStartUp() error {
	deleted, err := a.Config.Elastic.Index().Delete(a.Config.Index)
	glg.Infof("deleted index: %s %v", a.Config.Index, deleted)
	if err != nil {
		return err
	}

	return nil
}

func (a *AnaximanderHandler) CreateIndexAtStartup() error {
	indexMapping := a.Config.Elastic.Builder().Index()
	created, err := a.Config.Elastic.Index().Create(a.Config.Index, indexMapping)
	if err != nil {
		return err
	}

	glg.Infof("created index: %s %v", a.Config.Index, created.Acknowledged)

	return nil
}

func (a *AnaximanderHandler) AddToElastic(declension models.Declension, wg *sync.WaitGroup) error {
	defer wg.Done()
	upload, _ := declension.Marshal()

	_, err := a.Config.Elastic.Index().CreateDocument(a.Config.Index, upload)
	a.Config.Created++
	if err != nil {
		return err
	}

	return nil
}
