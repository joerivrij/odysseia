package models

import "encoding/json"

func UnmarshalWord(data []byte) (Word, error) {
	var r Word
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Word) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Word struct {
	Method      string `json:"method"`
	Category    string `json:"category"`
	Greek       string `json:"greek"`
	Translation string `json:"translation"`
	Chapter     int64  `json:"chapter"`
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
