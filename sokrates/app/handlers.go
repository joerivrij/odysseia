package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"net/http"
	"strings"
)

type SokratesHandler struct {
	Config *configs.SokratesConfig
}

// PingPong pongs the ping
func (s *SokratesHandler) PingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (s *SokratesHandler) health(w http.ResponseWriter, req *http.Request) {
	health := helpers.GetHealthOfApp(s.Config.ElasticClient)
	if !health.Healthy {
		middleware.ResponseWithCustomCode(w, 502, health)
		return
	}

	middleware.ResponseWithJson(w, health)
}

func (s *SokratesHandler) FindHighestChapter(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	category := pathParams["category"]
	method := pathParams["method"]

	if len(category) < 2 {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "category",
					Message: "must be longer than 1",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"method": method,
						},
					},
					{
						"match": map[string]interface{}{
							"category": category,
						},
					},
				},
			},
		},
	}

	elasticResult, err := elastic.QueryWithDescendingSort(s.Config.ElasticClient, s.Config.Index, "chapter", 1, query)
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

	elasticResult, err := elastic.QueryWithMatch(s.Config.ElasticClient, s.Config.Index, s.Config.SearchWord, checkAnswerRequest.QuizWord)
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
		if logos.Dutch == checkAnswerRequest.AnswerProvided {
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

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"method": method,
						},
					},
					{
						"match": map[string]interface{}{
							"category": category,
						},
					},
					{
						"match": map[string]interface{}{
							"chapter": chapter,
						},
					},
				},
			},
		},
	}

	elasticResponse, err := elastic.QueryWithScroll(s.Config.ElasticClient, s.Config.Index, query)
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
	for _, hit := range elasticResponse.Hits.Hits {
		source, _ := json.Marshal(hit.Source)
		logos, _ := models.UnmarshalWord(source)
		logoi.Logos = append(logoi.Logos, logos)
	}
	randNumber := helpers.GenerateRandomNumber(len(logoi.Logos))

	glg.Debugf("randomNumber: %d", randNumber)

	question := logoi.Logos[randNumber]
	quiz = append(quiz, question.Greek)
	quiz = append(quiz, question.Dutch)

	numberOfNeededAnswers := 5

	if len(logoi.Logos) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(logoi.Logos) + 1
	}

	for len(quiz) != numberOfNeededAnswers {
		randNumber = helpers.GenerateRandomNumber(len(logoi.Logos))
		randEntry := logoi.Logos[randNumber]

		exists := findQuizWord(quiz, randEntry.Dutch)
		if !exists {
			quiz = append(quiz, randEntry.Dutch)
		}
	}

	middleware.ResponseWithJson(w, quiz)
}

func (s *SokratesHandler) queryMethods(w http.ResponseWriter, req *http.Request) {
	field := "method.keyword"
	elasticResult, err := elastic.QueryUniqueField(s.Config.ElasticClient, field, s.Config.Index)

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
	method := pathParams["method"]
	field := "category.keyword"

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"method": method,
			},
		},
		"size": 0,
		"aggs": map[string]interface{}{
			"categories": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": field,
					"size":  500,
				},
			},
		},
	}

	elasticResult, err := elastic.QueryUniqueWithSupplied(s.Config.ElasticClient, s.Config.Index, query)

	glg.Info(elasticResult)
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
