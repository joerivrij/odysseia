package models

import "encoding/json"

func (r *CheckAnswerRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CheckAnswerRequest struct {
	QuizWord       string `json:"quizWord"`
	AnswerProvided string `json:"answerProvided"`
	Category       string `json:"category"`
}

type CheckAnswerResponse struct {
	Correct bool `json:"correct"`
}

type LastChapterResponse struct {
	LastChapter int64 `json:"lastChapter"`
}

type QuizResponse []string
