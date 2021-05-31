package impl

import (
	"encoding/json"
	"github.com/kpango/glg"
	"github.com/lexiko/alexandros/pkg/config"
	"github.com/lexiko/plato/elastic"
	"github.com/lexiko/plato/middleware"
	"github.com/lexiko/plato/models"
	"net/http"
)

type AlexandrosHandler struct {
	Config *config.AlexandrosConfig
}

// PingPong pongs the ping
func (a *AlexandrosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// Search a word based on part of that word
func (a *AlexandrosHandler) searchWord(w http.ResponseWriter, req *http.Request) {
	queryWord := req.URL.Query().Get("word")

	var searchResults []models.Meros

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

	glg.Debugf("looking for %s", queryWord)

	response, _ := elastic.QueryMultiMatchWithGrams(a.Config.ElasticClient, a.Config.Index, queryWord)

	if len(response.Hits.Hits) == 0 {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.NotFoundMessage{
				Type:   queryWord,
				Reason: "produced 0 results",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if meros.Original != "" {
			meros.Greek = meros.Original
			meros.Original = ""
		}
		searchResults = append(searchResults, meros)
	}

	middleware.ResponseWithJson(w, searchResults)
	return
}
