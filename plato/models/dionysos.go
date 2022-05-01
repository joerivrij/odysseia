package models

import "encoding/json"

type DeclensionConfig struct {
	Declensions []Declension
}

func UnmarshalDeclension(data []byte) (Declension, error) {
	var r Declension
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Declension) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Declension struct {
	Name        string              `json:"name"`
	Type        string              `json:"type,omitempty"`
	Dialect     string              `json:"dialect"`
	Declensions []DeclensionElement `json:"declensions"`
}

type DeclensionElement struct {
	Declension string   `json:"declension"`
	RuleName   string   `json:"ruleName"`
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
	Rule        string   `json:"rule,omitempty"`
	SearchTerms []string `json:"searchTerm,omitempty"`
}

func UnmarshalDeclensionTranslationResults(data []byte) (DeclensionTranslationResults, error) {
	var r DeclensionTranslationResults
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeclensionTranslationResults) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DeclensionTranslationResults struct {
	Results []Result `json:"results"`
}

type Result struct {
	Word        string `json:"word"`
	Rule        string `json:"rule"`
	RootWord    string `json:"rootWord"`
	Translation string `json:"translation"`
}

func (r *DeclensionTranslationResults) RemoveIndex(index int) {
	r.Results = append(r.Results[:index], r.Results[index+1:]...)
}
