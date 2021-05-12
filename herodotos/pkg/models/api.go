package models

type CreateSentenceResponse struct {
	Sentence   string `json:"sentence"`
	SentenceId string `json:"sentenceId"`
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
	SplitQuizSentence     []string           `json:"splitQuizSentence"`
	SplitAnswerSentence   []string           `json:"splitAnswerSentence"`
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
