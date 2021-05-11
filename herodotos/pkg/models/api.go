package models

type CreateQuestionResponse struct {
	Sentence   string `json:"sentence"`
	SentenceId string `json:"sentenceId"`
}

type CheckSentenceRequest struct {
	SentenceId       string `json:"sentenceId"`
	ProvidedSentence string `json:"answerSentence"`
	Author string `json:"author"`
}

type CheckSentenceResponse struct {
	LevenshteinPercentage string            `json:"levenshteinPercentage"`
	QuizSentence          string            `json:"quizSentence"`
	AnswerSentence        string            `json:"answerSentence"`
	MatchingWords         []MatchingWord    `json:"matchingWords,omitempty"`
	NonMatchingWords      []NonMatchingWord `json:"nonMatchingWords,omitempty"`
}

type MatchingWord struct {
	Word  string `json:"word"`
	Index int64  `json:"index"`
}

type NonMatchingWord struct {
	Word    string  `json:"word"`
	Index   int64   `json:"index"`
	Matches []Match `json:"matches"`
}

type Match struct {
	Match       string `json:"match"`
	Levenshtein int64  `json:"levenshtein"`
	Index       int64  `json:"index"`
	Percentage  float64 `json:"percentage"`
}

