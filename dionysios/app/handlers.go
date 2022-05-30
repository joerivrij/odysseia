package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"net/http"
)

type DionysosHandler struct {
	Config *configs.DionysiosConfig
}

// PingPong pongs the ping
func (d *DionysosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (d *DionysosHandler) health(w http.ResponseWriter, req *http.Request) {
	health := helpers.GetHealthOfApp(d.Config.Elastic)
	if !health.Healthy {
		middleware.ResponseWithCustomCode(w, 502, health)
		return
	}

	middleware.ResponseWithJson(w, health)
}

func (d *DionysosHandler) checkGrammar(w http.ResponseWriter, req *http.Request) {
	queryWord := req.URL.Query().Get("word")

	if queryWord == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "word",
					Message: "cannot be empty",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	glg.Debugf("trying to get the possibilities for %s", queryWord)

	cacheItem, _ := d.Config.Cache.Read(queryWord)
	if cacheItem != nil {
		d, err := models.UnmarshalDeclensionTranslationResults(cacheItem)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
				Messages: []models.ValidationMessages{
					{
						Field:   "cache",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithJson(w, e)
			return
		}
		middleware.ResponseWithJson(w, d)
		return
	}

	declensions, _ := d.StartFindingRules(queryWord)
	if len(declensions.Results) == 0 || declensions.Results == nil {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.NotFoundMessage{
				Type:   queryWord,
				Reason: "no options found",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	stringifiedDeclension, _ := declensions.Marshal()
	err := d.Config.Cache.Set(queryWord, string(stringifiedDeclension))
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "cache",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	middleware.ResponseWithJson(w, *declensions)
}
