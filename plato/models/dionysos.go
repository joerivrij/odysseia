package models

import "encoding/json"

func UnmarshalDeclensions(data []byte) (Declensions, error) {
	var r Declensions
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Declensions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Declensions struct {
	Declensions Declension `json:"declensions"`
}

type Declension struct {
	FirstDeclension  []DeclensionElement `json:"firstDeclension,omitempty"`
	SecondDeclension []DeclensionElement `json:"secondDeclension,omitempty"`
}

type DeclensionElement struct {
	Declension string       `json:"declension"`
	RuleName   string       `json:"ruleName"`
	SearchTerm []string `json:"searchTerm"`
}

func UnmarshalFoundRules(data []byte) (FoundRules, error) {
	var r FoundRules
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FoundRules) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FoundRules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Form       string `json:"form,omitempty"`
	Declension string `json:"declension,omitempty"`
	Rule       string `json:"rule,omitempty"`
	SearchTerms []string `json:"searchTerm,omitempty"`
}