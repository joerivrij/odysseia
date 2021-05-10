package models

type CreateQuestionResponse struct {
	Sentence       string `json:"sentence"`
	SentenceId	   string `json:"sentenceId"`
}

type CheckSentenceRequest struct {
	SentenceId	   string `json:"sentenceId"`
	ProvidedSentence string `json:"answerSentence"`
	QuizSentence string `json:"quizSentence"`
}

type CheckSentenceResponse struct {
	LevenshteinDistance       int `json:"levenshteinDistance"`
	LevenshteinPercentage       string `json:"levenshteinPercentage"`
	QuizSentence	   string `json:"quizSentence"`
	AnswerSentence	   string `json:"answerSentence"`
	MatchingWords []string `json:"matchingWords,omitempty"`
	MatchingWordsWithTypoAllowance []string `json:"matchingWordsWithTypoAllowance,omitempty"`
}