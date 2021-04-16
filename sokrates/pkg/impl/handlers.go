package impl

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"net/http"
	"sokrates/pkg/config"
	"sokrates/pkg/middleware"
	"sokrates/pkg/models"
)

type SokratesHandler struct {
	Config *config.SokratesConfig
}

// PingPong pongs the ping
func (s *SokratesHandler)PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (s *SokratesHandler)FindHighestChapter(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	category := pathParams["category"]

	chapter, elasticErr := QueryLastChapter(s.Config.ElasticClient, category)

	if elasticErr != nil {

		notFoundErr := models.NotFoundMessage{
			Type:   fmt.Sprintf("%v", elasticErr["error"].(map[string]interface{})["type"]),
			Reason: fmt.Sprintf("%v", elasticErr["error"].(map[string]interface{})["reason"]),
		}
		e := models.NotFoundError{ErrorModel: models.ErrorModel{middleware.CreateGUID()}, Message: notFoundErr}
		middleware.ResponseWithJson(w, e)
		return
	}

	response := models.LastChapterResponse{LastChapter: chapter}

	middleware.ResponseWithJson(w, response)
}

func (s *SokratesHandler)CheckAnswer(w http.ResponseWriter, req *http.Request) {
	var checkAnswerRequest models.CheckAnswerRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkAnswerRequest)
	if err != nil {
		glg.Error(err)
	}

	storedAnswer, _ := QueryWithMatch(s.Config.ElasticClient, checkAnswerRequest.Category, s.Config.SearchTerm, checkAnswerRequest.QuizWord)

	answer := models.CheckAnswerResponse{ Correct : false}

	for _, logos := range storedAnswer.Logos {
		if logos.Dutch == checkAnswerRequest.AnswerProvided {
			answer.Correct = true
		}
	}

	middleware.ResponseWithJson(w, answer)
}

func (s *SokratesHandler)CreateQuestion(w http.ResponseWriter, req *http.Request) {
	chapter := req.URL.Query().Get("chapter")
	category := req.URL.Query().Get("category")

	questionSet, _ := QueryWithScroll(s.Config.ElasticClient, category, "chapter", chapter)

	fmt.Print(questionSet)

	middleware.ResponseWithJson(w, chapter)
}
