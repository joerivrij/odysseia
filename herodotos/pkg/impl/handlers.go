package impl

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/lexiko/herodotos/pkg/config"
	apiModels "github.com/lexiko/herodotos/pkg/models"
	"github.com/lexiko/plato/elastic"
	"github.com/lexiko/plato/helpers"
	"github.com/lexiko/plato/middleware"
	"github.com/lexiko/plato/models"
	"net/http"
	"strings"
)

type HerodotosHandler struct {
	Config *config.HerodotosConfig
}

// PingPong pongs the ping
func (h *HerodotosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// creates a new sentence for questions
func (h *HerodotosHandler) createQuestion(w http.ResponseWriter, req *http.Request) {
	author := req.URL.Query().Get("author")

	if author == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "author",
					Message: "cannot be empty",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	response, _ := elastic.QueryWithMatchAll(h.Config.ElasticClient, author)
	randNumber := helpers.GenerateRandomNumber(len(response.Hits.Hits))
	questionItem := response.Hits.Hits[randNumber]
	id := questionItem.ID

	elasticJson, _ := json.Marshal(questionItem.Source)
	rhemaSource, err := models.UnmarshalRhema(elasticJson)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "",
					Message: "something went wrong",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}
	glg.Info(rhemaSource.Greek)

	question := apiModels.CreateSentenceResponse{Sentence: rhemaSource.Greek,
		SentenceId: id}
	middleware.ResponseWithJson(w, question)
}

// checks the validity of an answer
func (h *HerodotosHandler) checkSentence(w http.ResponseWriter, req *http.Request) {
	var checkSentenceRequest apiModels.CheckSentenceRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkSentenceRequest)
	if err != nil {
		glg.Error(err)
	}

	elasticResult, err := elastic.QueryOnId(h.Config.ElasticClient, strings.ToLower(checkSentenceRequest.Author), checkSentenceRequest.SentenceId)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "s",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	elasticJson, _ := json.Marshal(elasticResult.Hits.Hits[0].Source)
	original, err := models.UnmarshalRhema(elasticJson)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages: []models.ValidationMessages{
				{
					Field:   "s",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

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

	response := apiModels.CheckSentenceResponse{
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
	elasticResult, _ := elastic.QueryWithMatchAll(h.Config.ElasticClient, h.Config.AuthorIndex)
	var authors models.Authors
	for _, hit := range elasticResult.Hits.Hits {
		elasticJson, _ := json.Marshal(hit.Source)
		author, err := models.UnmarshalAuthors(elasticJson)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
				Messages: []models.ValidationMessages{
					{
						Field:   "s",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithJson(w, e)
			return
		}
		authors.Authors = append(authors.Authors, author)
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
func findMatchingWordsWithSpellingAllowance(source, target string) (response apiModels.CheckSentenceResponse) {
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
				response.MatchingWords = append(response.MatchingWords, apiModels.MatchingWord{
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
			levenshteinModel := apiModels.NonMatchingWord{
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

				matchModel := apiModels.Match{
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
		if strings.IndexRune(characters, r) < 0 {
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
