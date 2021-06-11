package impl

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"github.com/lexiko/plato/helpers"
	"github.com/lexiko/plato/models"
	"github.com/lexiko/sokrates/pkg/config"
	"github.com/lexiko/sokrates/pkg/middleware"
	apiModels "github.com/lexiko/sokrates/pkg/models"
	"net/http"
)

type SokratesHandler struct {
	Config *config.SokratesConfig
}

// PingPong pongs the ping
func (s *SokratesHandler) PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (s *SokratesHandler) FindHighestChapter(w http.ResponseWriter, req *http.Request) {
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

	response := apiModels.LastChapterResponse{LastChapter: chapter}

	middleware.ResponseWithJson(w, response)
}

func (s *SokratesHandler) CheckAnswer(w http.ResponseWriter, req *http.Request) {
	var checkAnswerRequest apiModels.CheckAnswerRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkAnswerRequest)
	if err != nil {
		glg.Error(err)
	}

	storedAnswer, _ := QueryWithMatch(s.Config.ElasticClient, checkAnswerRequest.Category, s.Config.SearchTerm, checkAnswerRequest.QuizWord)

	answer := apiModels.CheckAnswerResponse{Correct: false}

	for _, logos := range storedAnswer.Logos {
		if logos.Dutch == checkAnswerRequest.AnswerProvided {
			answer.Correct = true
		}
	}

	middleware.ResponseWithJson(w, answer)
}

func (s *SokratesHandler) CreateQuestion(w http.ResponseWriter, req *http.Request) {
	chapter := req.URL.Query().Get("chapter")
	category := req.URL.Query().Get("category")

	glg.Debugf("category: %s chapter: %s", category, chapter)

	var quiz apiModels.QuizResponse

	questionSet, _ := QueryWithScroll(s.Config.ElasticClient, category, "chapter", chapter)
	randNumber := helpers.GenerateRandomNumber(len(questionSet.Logos))

	glg.Debugf("randomNumber: %d", randNumber)

	question := questionSet.Logos[randNumber]
	quiz = append(quiz, question.Greek)
	quiz = append(quiz, question.Dutch)

	numberOfNeededAnswers := 5

	if len(questionSet.Logos) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(questionSet.Logos) + 1
	}

	for len(quiz) != numberOfNeededAnswers {
		randNumber = helpers.GenerateRandomNumber(len(questionSet.Logos))
		randEntry := questionSet.Logos[randNumber]

		exists := findQuizWord(quiz, randEntry.Dutch)
		if !exists {
			quiz = append(quiz, randEntry.Dutch)
		}
	}

	middleware.ResponseWithJson(w, quiz)
}

// findQuizWord takes a slice and looks for an element in it
func findQuizWord(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
