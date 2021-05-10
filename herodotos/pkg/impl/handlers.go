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
	response, _ := elastic.QueryWithMatchAll(h.Config.ElasticClient, h.Config.ElasticIndex)
	randNumber := helpers.GenerateRandomNumber(len(response.Hits.Hits))
	questionItem := response.Hits.Hits[randNumber]
	id := questionItem.ID

	elasticJson, _ := json.Marshal(questionItem.Source)
	rhemaSource, err := models.UnmarshalRhema(elasticJson)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Messages:   []models.ValidationMessages{
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

	question := apiModels.CreateQuestionResponse{Sentence: rhemaSource.Greek,
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

	matchingWords, matchingWordsWithAllowance := findMatchingWordsWithSpellingAllowance(checkSentenceRequest.QuizSentence, checkSentenceRequest.ProvidedSentence)

	levenshteinDist := levenshteinDistance(checkSentenceRequest.QuizSentence, checkSentenceRequest.ProvidedSentence)
	lenOfLongestSentence := longestStringOfTwo(checkSentenceRequest.QuizSentence, checkSentenceRequest.ProvidedSentence)
	percentage := levenshteinDistanceInPercentage(levenshteinDist, lenOfLongestSentence)
	roundedPercentage := fmt.Sprintf("%.2f", percentage)
	glg.Infof("levenshtein percentage: %s", roundedPercentage)

	response := apiModels.CheckSentenceResponse{
		LevenshteinDistance:   levenshteinDist,
		LevenshteinPercentage: roundedPercentage,
		QuizSentence:              checkSentenceRequest.QuizSentence,
		AnswerSentence: checkSentenceRequest.ProvidedSentence,
		MatchingWords: matchingWords,
		MatchingWordsWithTypoAllowance: matchingWordsWithAllowance,
	}
	middleware.ResponseWithJson(w, response)
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
func levenshteinDistanceInPercentage(levenshteinDistance int, longestString int) float32 {
	return (1.00 - float32(levenshteinDistance)/float32(longestString)) * 100.00
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
func findMatchingWordsWithSpellingAllowance(source, target string) (matchingWords []string, matchingWordsWithAllowance[]string) {
	sourceSentence := strings.Split(source, " ")
	targetSentence := strings.Split(target, " ")

	for _, wordInSource := range sourceSentence {
		for _, wordInTarget := range targetSentence {
			sourceWord := strings.ToLower(wordInSource)
			targetWord := strings.ToLower(wordInTarget)
			wordInSlice := checkSliceForItem(matchingWords, wordInSource)
			wordInLevenshteinSlice := checkSliceForItem(matchingWordsWithAllowance, wordInSource)
			if wordInSlice {
				continue
			} else if sourceWord == targetWord {
				matchingWords = append(matchingWords, wordInSource)
				break
			}
			if wordInLevenshteinSlice {
				continue
			}
			levenshtein := levenshteinDistance(sourceWord, targetWord)
			if levenshtein == 1 {
				matchingWordsWithAllowance = append(matchingWordsWithAllowance, wordInTarget)
			}
		}
	}

	return matchingWords, matchingWordsWithAllowance
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

func streamlineSentenceBeforeCompare(matchingWords []string, sentence string) (string) {
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