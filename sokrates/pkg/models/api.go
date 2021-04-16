package models

import "encoding/json"

type ResultModel struct {
	Result string `json:"result"`
}

func UnmarshalWord(data []byte) (Word, error) {
	var r Word
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Word) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Word struct {
	Greek   string `json:"greek"`
	Dutch   string `json:"dutch"`
	Chapter int64  `json:"chapter"`
}

func UnmarshalLogos(data []byte) (Logos, error) {
	var r Logos
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Logos) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Logos struct {
	Logos []Word `json:"logos"`
}

func UnmarshalCheckAnswerRequest(data []byte) (CheckAnswerRequest, error) {
	var r CheckAnswerRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CheckAnswerRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CheckAnswerRequest struct {
	QuizWord       string `json:"quizWord"`
	AnswerProvided string `json:"answerProvided"`
	Category      string `json:"category"`
}

type CheckAnswerResponse struct {
	Correct bool `json:"correct"`
}

type LastChapterResponse struct {
	LastChapter int64 `json:"lastChapter"`
}
