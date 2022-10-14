package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/helpers"
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia/aristoteles/configs"
	"net/http"
	"strings"
)

type SokratesHandler struct {
	Config *configs.SokratesConfig
}

const (
	Method     string = "method"
	Authors    string = "authors"
	Category   string = "category"
	Categories string = "categories"
	Chapter    string = "chapter"
)

// PingPong pongs the ping
func (s *SokratesHandler) PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (s *SokratesHandler) health(w http.ResponseWriter, req *http.Request) {
	health := helpers.GetHealthOfApp(s.Config.Elastic)
	if !health.Healthy {
		middleware.ResponseWithCustomCode(w, 502, health)
		return
	}

	middleware.ResponseWithJson(w, health)
}

func (s *SokratesHandler) FindHighestChapter(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	category := pathParams[Category]
	method := pathParams[Method]

	if len(category) < 2 {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   Category,
					Message: "must be longer than 1",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	mustQuery := []map[string]string{
		{
			Method: method,
		},
		{
			Category: category,
		},
	}

	query := s.Config.Elastic.Builder().MultipleMatch(mustQuery)
	mode := "desc"

	elasticResult, err := s.Config.Elastic.Query().MatchWithSort(s.Config.Index, mode, Chapter, 1, query)
	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	elasticJson, _ := json.Marshal(elasticResult.Hits.Hits[0].Source)
	chapter, _ := models.UnmarshalWord(elasticJson)
	response := models.LastChapterResponse{LastChapter: chapter.Chapter}

	middleware.ResponseWithJson(w, response)
}

func (s *SokratesHandler) CheckAnswer(w http.ResponseWriter, req *http.Request) {
	var checkAnswerRequest models.CheckAnswerRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkAnswerRequest)
	if err != nil || checkAnswerRequest.AnswerProvided == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "body",
					Message: "error parsing",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	query := s.Config.Elastic.Builder().MatchQuery(s.Config.SearchWord, checkAnswerRequest.QuizWord)
	elasticResult, err := s.Config.Elastic.Query().Match(s.Config.Index, query)
	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}
	var logoi models.Logos
	answer := models.CheckAnswerResponse{Correct: false}
	for _, hit := range elasticResult.Hits.Hits {
		elasticJson, _ := json.Marshal(hit.Source)
		logos, _ := models.UnmarshalWord(elasticJson)
		logoi.Logos = append(logoi.Logos, logos)
	}

	for _, logos := range logoi.Logos {
		if logos.Translation == checkAnswerRequest.AnswerProvided {
			answer.Correct = true
		}
	}

	middleware.ResponseWithJson(w, answer)
}

func (s *SokratesHandler) CreateQuestion(w http.ResponseWriter, req *http.Request) {
	chapter := req.URL.Query().Get("chapter")
	category := req.URL.Query().Get("category")
	method := req.URL.Query().Get("method")

	if category == "" || chapter == "" || method == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "category, chapter, method",
					Message: "cannot be empty",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	glg.Debugf("category: %s chapter: %s method: %s", category, chapter, method)

	var quiz models.QuizResponse

	mustQuery := []map[string]string{
		{
			Method: method,
		},
		{
			Category: category,
		},
		{
			Chapter: chapter,
		},
	}

	query := s.Config.Elastic.Builder().MultipleMatch(mustQuery)

	elasticResponse, err := s.Config.Elastic.Query().MatchWithScroll(s.Config.Index, query)

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			e := models.NotFoundError{
				ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
				Message: models.NotFoundMessage{
					Type:   "no results",
					Reason: fmt.Sprintf("category: %s chapter: %s method: %s", category, chapter, method),
				},
			}
			middleware.ResponseWithJson(w, e)
			return
		}
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var logoi models.Logos
	for _, hit := range elasticResponse.Hits.Hits {
		source, _ := json.Marshal(hit.Source)
		logos, _ := models.UnmarshalWord(source)
		logoi.Logos = append(logoi.Logos, logos)
	}
	randNumber := helpers.GenerateRandomNumber(len(logoi.Logos))

	glg.Debugf("randomNumber: %d", randNumber)

	question := logoi.Logos[randNumber]
	quiz = append(quiz, question.Greek)
	quiz = append(quiz, question.Translation)

	numberOfNeededAnswers := 5

	if len(logoi.Logos) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(logoi.Logos) + 1
	}

	for len(quiz) != numberOfNeededAnswers {
		randNumber = helpers.GenerateRandomNumber(len(logoi.Logos))
		randEntry := logoi.Logos[randNumber]

		exists := findQuizWord(quiz, randEntry.Translation)
		if !exists {
			quiz = append(quiz, randEntry.Translation)
		}
	}

	middleware.ResponseWithJson(w, quiz)
}

func (s *SokratesHandler) queryMethods(w http.ResponseWriter, req *http.Request) {
	field := "method.keyword"
	query := s.Config.Elastic.Builder().Aggregate(Authors, field)
	elasticResult, err := s.Config.Elastic.Query().MatchAggregate(s.Config.Index, query)

	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var methods models.Methods
	for _, bucket := range elasticResult.Aggregations.AuthorAggregation.Buckets {
		author := models.Method{Method: strings.ToLower(fmt.Sprintf("%v", bucket.Key))}
		methods.Method = append(methods.Method, author)
	}

	middleware.ResponseWithJson(w, methods)
}

func (s *SokratesHandler) queryCategories(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	method := pathParams[Method]
	field := fmt.Sprintf("%s.keyword", Category)

	query := s.Config.Elastic.Builder().FilteredAggregate(Method, method, Categories, field)
	elasticResult, err := s.Config.Elastic.Query().MatchAggregate(s.Config.Index, query)

	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var categories models.Categories
	for _, bucket := range elasticResult.Aggregations.CategoryAggregation.Buckets {
		category := models.Category{Category: fmt.Sprintf("%s", bucket.Key)}
		categories.Category = append(categories.Category, category)

	}

	middleware.ResponseWithJson(w, categories)
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
