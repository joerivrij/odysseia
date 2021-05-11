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
	Author       string   `json:"author"`
	Greek        string   `json:"greek"`
	Translations []string `json:"translations"`
}

type Rhema struct {
	Rhemai []RhemaSource `json:"rhemai"`
}
