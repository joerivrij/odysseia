package models

import "encoding/json"

func (r *Biblos) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Meros) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalBiblos(data []byte) (Biblos, error) {
	var r Biblos
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalMeros(data []byte) (Meros, error) {
	var r Meros
	err := json.Unmarshal(data, &r)
	return r, err
}

type Biblos struct {
	Biblos []Meros `json:"biblos"`
}

type Meros struct {
	Greek      string `json:"greek"`
	English    string `json:"english"`
	LinkedWord string `json:"linkedWord,omitempty"`
	Original   string `json:"original,omitempty"`
}
