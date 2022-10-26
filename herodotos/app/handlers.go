package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/helpers"
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"net/http"
	"strings"
)

type HerodotosHandler struct {
	Config *configs.HerodotosConfig
}

const (
	Author  string = "author"
	Authors string = "authors"
	Book    string = "book"
	Books   string = "books"
	Id      string = "_id"
)

// PingPong pongs the ping
func (h *HerodotosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (h *HerodotosHandler) health(w http.ResponseWriter, req *http.Request) {
	health := helpers.GetHealthOfApp(h.Config.Elastic)
	if !health.Healthy {
		middleware.ResponseWithCustomCode(w, 502, health)
		return
	}

	middleware.ResponseWithJson(w, health)
}

// creates a new sentence for questions
func (h *HerodotosHandler) createQuestion(w http.ResponseWriter, req *http.Request) {
	author := req.URL.Query().Get(Author)
	book := req.URL.Query().Get(Book)

	if author == "" || book == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "author and book",
					Message: "cannot be empty",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	mustQuery := []map[string]string{
		{
			Author: author,
		},
		{
			Book: book,
		},
	}

	query := h.Config.Elastic.Builder().MultipleMatch(mustQuery)

	response, err := h.Config.Elastic.Query().Match(h.Config.Index, query)

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

	if len(response.Hits.Hits) == 0 {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.NotFoundMessage{
				Type:   fmt.Sprintf("author: %s and book: %s", author, book),
				Reason: "no hits for combination",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	randNumber := helpers.GenerateRandomNumber(len(response.Hits.Hits))
	questionItem := response.Hits.Hits[randNumber]
	id := questionItem.ID

	elasticJson, _ := json.Marshal(questionItem.Source)
	rhemaSource, err := models.UnmarshalRhema(elasticJson)
	if err != nil || rhemaSource.Translations == nil {
		errorMessage := fmt.Errorf("an error occurred while parsing %s", elasticJson)
		glg.Error(errorMessage.Error())
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "createQuestion",
					Message: errorMessage.Error(),
				},
				{
					Field:   "translation",
					Message: "cannot be nil",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	question := models.CreateSentenceResponse{Sentence: rhemaSource.Greek,
		SentenceId: id}
	middleware.ResponseWithJson(w, question)
}

// checks the validity of an answer
func (h *HerodotosHandler) checkSentence(w http.ResponseWriter, req *http.Request) {
	var checkSentenceRequest models.CheckSentenceRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkSentenceRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "decoding",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	query := h.Config.Elastic.Builder().MatchQuery(Id, checkSentenceRequest.SentenceId)
	elasticResult, err := h.Config.Elastic.Query().Match(h.Config.Index, query)
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
	original, err := models.UnmarshalRhema(elasticJson)
	if err != nil || original.Translations == nil {
		errorMessage := fmt.Errorf("an error occurred while parsing %s", elasticJson)
		glg.Error(errorMessage.Error())
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "createQuestion",
					Message: errorMessage.Error(),
				},
				{
					Field:   "translation",
					Message: "cannot be nil",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	glg.Info(original.Greek)

	var sentence string
	var percentage float64
	for _, solution := range original.Translations {
		levenshteinDist := levenshteinDistance(solution, checkSentenceRequest.ProvidedSentence)
		lenOfLongestSentence := longestStringOfTwo(solution, checkSentenceRequest.ProvidedSentence)
		levenshteinPerc := levenshteinDistanceInPercentage(levenshteinDist, lenOfLongestSentence)
		if levenshteinPerc > percentage {
			sentence = solution
			percentage = levenshteinPerc
		}
	}

	roundedPercentage := fmt.Sprintf("%.2f", percentage)
	glg.Infof("levenshtein percentage: %s", roundedPercentage)

	model := findMatchingWordsWithSpellingAllowance(sentence, checkSentenceRequest.ProvidedSentence)

	response := models.CheckSentenceResponse{
		LevenshteinPercentage: roundedPercentage,
		QuizSentence:          sentence,
		AnswerSentence:        checkSentenceRequest.ProvidedSentence,
		MatchingWords:         model.MatchingWords,
		NonMatchingWords:      model.NonMatchingWords,
		SplitQuizSentence:     model.SplitQuizSentence,
		SplitAnswerSentence:   model.SplitAnswerSentence,
	}

	middleware.ResponseWithJson(w, response)
}

func (h *HerodotosHandler) queryAuthors(w http.ResponseWriter, req *http.Request) {
	field := fmt.Sprintf("%s.keyword", Author)
	query := h.Config.Elastic.Builder().Aggregate(Authors, field)

	elasticResult, err := h.Config.Elastic.Query().MatchAggregate(h.Config.Index, query)

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

	var authors models.Authors
	for _, bucket := range elasticResult.Aggregations.AuthorAggregation.Buckets {
		author := models.Author{Author: strings.ToLower(fmt.Sprintf("%v", bucket.Key))}
		authors.Authors = append(authors.Authors, author)
	}

	middleware.ResponseWithJson(w, authors)
}

func (h *HerodotosHandler) queryBooks(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	author := pathParams[Author]

	query := h.Config.Elastic.Builder().FilteredAggregate(Author, author, Books, Book)
	elasticResult, err := h.Config.Elastic.Query().MatchAggregate(h.Config.Index, query)

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

	var authors models.Books
	for _, bucket := range elasticResult.Aggregations.BookAggregation.Buckets {

		book := models.Book{Book: int64(bucket.Key.(float64))}
		authors.Books = append(authors.Books, book)

	}

	middleware.ResponseWithJson(w, authors)
}

// calculates the amount of changes needed to have two sentences match
// example: Distance from Python to Peithen is 3
func levenshteinDistance(question, answer string) int {
	questionLen := len(question)
	answerLen := len(answer)
	column := make([]int, len(question)+1)

	for y := 1; y <= questionLen; y++ {
		column[y] = y
	}
	for x := 1; x <= answerLen; x++ {
		column[0] = x
		lastKey := x - 1
		for y := 1; y <= questionLen; y++ {
			oldKey := column[y]
			var incr int
			if question[y-1] != answer[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastKey+incr)
			lastKey = oldKey
		}
	}
	return column[questionLen]
}

// creates a percentage based on the levenshtein and the longest string
func levenshteinDistanceInPercentage(levenshteinDistance int, longestString int) float64 {
	return (1.00 - float64(levenshteinDistance)/float64(longestString)) * 100.00
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

func longestStringOfTwo(a, b string) int {
	if len(a) >= len(b) {
		return len(a)
	}
	return len(b)
}

// take two sentences and creates a list of matching words and words with a typo (1 levenshtein)
func findMatchingWordsWithSpellingAllowance(source, target string) (response models.CheckSentenceResponse) {
	s := removeCharacters(source, ",`~<>/?!.;:'\"")
	t := removeCharacters(target, ",`~<>/?!.;:'\"")

	sourceSentence := strings.Split(s, " ")
	targetSentence := strings.Split(t, " ")

	response.SplitQuizSentence = sourceSentence
	response.SplitAnswerSentence = targetSentence

	for i, wordInSource := range sourceSentence {
		for _, wordInTarget := range targetSentence {
			sourceWord := strings.ToLower(wordInSource)
			targetWord := strings.ToLower(wordInTarget)

			levenshtein := levenshteinDistance(sourceWord, targetWord)

			// might be changed to one levenshtein as being typo's
			if levenshtein == 0 {
				response.MatchingWords = append(response.MatchingWords, models.MatchingWord{
					Word:        wordInSource,
					SourceIndex: i,
				})
				break
			}
		}
	}

	var slice []string
	for _, word := range response.MatchingWords {
		slice = append(slice, word.Word)
	}

	for i, wordInSource := range sourceSentence {
		wordInSlice := checkSliceForItem(slice, wordInSource)
		if !wordInSlice {
			levenshteinModel := models.NonMatchingWord{
				Word:        wordInSource,
				SourceIndex: i,
				Matches:     nil,
			}
			for j, wordInTarget := range targetSentence {
				sourceWord := strings.ToLower(wordInSource)
				targetWord := strings.ToLower(wordInTarget)

				levenshtein := levenshteinDistance(sourceWord, targetWord)
				percentage := levenshteinDistanceInPercentage(levenshtein, longestStringOfTwo(sourceWord, targetWord))
				roundedPercentage := fmt.Sprintf("%.2f", percentage)

				matchModel := models.Match{
					Match:       wordInTarget,
					Levenshtein: levenshtein,
					AnswerIndex: j,
					Percentage:  roundedPercentage,
				}
				levenshteinModel.Matches = append(levenshteinModel.Matches, matchModel)
			}
			response.NonMatchingWords = append(response.NonMatchingWords, levenshteinModel)
		}
	}

	return
}

// takes a slice and returns a bool if value is part of the slice
func checkSliceForItem(slice []string, sourceWord string) bool {
	for _, item := range slice {
		if item == sourceWord {
			return true
		}
	}

	return false
}

func removeCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if !strings.ContainsRune(characters, r) {
			return r
		}
		return -1
	}

	return strings.Map(filter, input)

}

func streamlineSentenceBeforeCompare(matchingWords []string, sentence string) string {
	newSentence := ""

	sourceSentence := strings.Split(sentence, " ")

	for index, wordInSource := range sourceSentence {
		wordInSlice := checkSliceForItem(matchingWords, wordInSource)
		if !wordInSlice {
			newSentence += wordInSource
			if len(sourceSentence) != index+1 {
				newSentence += " "
			}
		}
	}

	return newSentence
}
