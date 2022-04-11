package models

import "encoding/json"

func UnmarshalRhema(data []byte) (RhemaSource, error) {
	var r RhemaSource
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RhemaSource) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RhemaSource struct {
	Author          string   `json:"author"`
	Greek           string   `json:"greek"`
	Translations    []string `json:"translations"`
	Book            int64    `json:"book"`
	Chapter         int64    `json:"chapter"`
	Section         int64    `json:"section"`
	PerseusTextLink string   `json:"perseusTextLink"`
}

type Rhema struct {
	Rhemai []RhemaSource `json:"rhemai"`
}

func UnmarshalAuthors(data []byte) (Author, error) {
	var r Author
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Author) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Authors struct {
	Authors []Author `json:"authors"`
}

type Author struct {
	Author string `json:"author"`
}

type Books struct {
	Books []Book `json:"books"`
}

type Book struct {
	Book int64 `json:"book"`
}

type Methods struct {
	Method []Method `json:"methods"`
}

type Method struct {
	Method string `json:"method"`
}

type Categories struct {
	Category []Category `json:"categories"`
}

type Category struct {
	Category string `json:"category"`
}

type CreateSentenceResponse struct {
	Sentence   string `json:"sentence"`
	SentenceId string `json:"sentenceId"`
}

func (r *CheckSentenceRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CheckSentenceRequest struct {
	SentenceId       string `json:"sentenceId"`
	ProvidedSentence string `json:"answerSentence"`
	Author           string `json:"author"`
}

type CheckSentenceResponse struct {
	LevenshteinPercentage string            `json:"levenshteinPercentage"`
	QuizSentence          string            `json:"quizSentence"`
	AnswerSentence        string            `json:"answerSentence"`
	SplitQuizSentence     []string          `json:"splitQuizSentence"`
	SplitAnswerSentence   []string          `json:"splitAnswerSentence"`
	MatchingWords         []MatchingWord    `json:"matchingWords,omitempty"`
	NonMatchingWords      []NonMatchingWord `json:"nonMatchingWords,omitempty"`
}

type MatchingWord struct {
	Word        string `json:"word"`
	SourceIndex int    `json:"sourceIndex"`
}

type NonMatchingWord struct {
	Word        string  `json:"word"`
	SourceIndex int     `json:"sourceIndex"`
	Matches     []Match `json:"matches"`
}

type Match struct {
	Match       string `json:"match"`
	Levenshtein int    `json:"levenshtein"`
	AnswerIndex int    `json:"answerIndex"`
	Percentage  string `json:"percentage"`
}
