package app

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/middleware"
	"github.com/odysseia/plato/models"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
	"unicode/utf8"
)

func (d *DionysosHandler) queryWordInElastic(word string) ([]models.Meros, error) {
	var searchResults []models.Meros
	strippedWord := d.removeAccents(word)

	term := "greek"
	query := d.Config.Elastic.Builder().MatchQuery(term, strippedWord)
	response, err := d.Config.Elastic.Query().Match(d.Config.SecondaryIndex, query)

	if err != nil {
		errText := err.Error()
		//todo better way to check for a 404
		if strings.Contains(errText, "404") {
			e := models.NotFoundError{
				ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
				Message: models.NotFoundMessage{
					Type:   term,
					Reason: errText,
				},
			}
			return nil, &e
		}

		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: middleware.CreateGUID()},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		return nil, &e
	}

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if meros.Original != "" {
			meros.Greek = meros.Original
			meros.Original = ""
		}
		searchResults = append(searchResults, meros)
	}

	return searchResults, nil
}

func (d *DionysosHandler) removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		glg.Error(e.Error())
	}
	return output
}

func (d *DionysosHandler) parseDictResults(dictionaryHits models.Meros) (translation, article string) {
	translation = dictionaryHits.English
	greek := strings.Split(dictionaryHits.Greek, ",")
	if len(greek) > 1 {
		article = strings.Replace(greek[1], " ", "", -1)
	}

	return
}

func (d *DionysosHandler) StartFindingRules(word string) (*models.DeclensionTranslationResults, error) {
	var results models.DeclensionTranslationResults

	declensions, err := d.searchForDeclensions(word)
	if err != nil {
		return nil, err
	}

	if len(declensions.Rules) > 0 {
		for _, declension := range declensions.Rules {
			if len(declension.SearchTerms) > 0 {
				for _, term := range declension.SearchTerms {
					dictionaryHits, err := d.queryWordInElastic(term)
					if err != nil {
						glg.Debug("word not found in database")
					}
					for _, hit := range dictionaryHits {
						translation, article := d.parseDictResults(hit)

						result := models.Result{
							Word:        word,
							Rule:        declension.Rule,
							RootWord:    term,
							Translation: translation,
						}

						if len(results.Results) > 0 {
							if translation == "" {
								continue
							}
						}

						if article != "" {
							switch article {
							case "ὁ":
								if !strings.Contains(declension.Rule, "masc") {
									continue
								}
							case "ἡ":
								if !strings.Contains(declension.Rule, "fem") {
									continue
								}
							case "τό":
								if !strings.Contains(declension.Rule, "neut") {
									continue
								}
							default:
								continue
							}
						}
						results.Results = append(results.Results, result)
					}

				}
			}
		}
	} else {
		singleSearchResult, err := d.queryWordInElastic(word)
		if err != nil {
			glg.Debug("no result for single word continuing loop")
		}

		if len(singleSearchResult) > 0 {
			for _, searchResult := range singleSearchResult {
				translation, _ := d.parseDictResults(searchResult)
				doNotAdd := false
				for _, res := range results.Results {
					if res.Translation == translation {
						doNotAdd = true
						break
					}
				}

				if doNotAdd {
					continue
				}

				result := models.Result{
					Word:        word,
					Rule:        "preposition",
					RootWord:    searchResult.Greek,
					Translation: translation,
				}
				results.Results = append(results.Results, result)
			}
		}
	}

	if results.Results == nil {
		singleSearchResult, err := d.queryWordInElastic(word)
		if err != nil {
			glg.Debug("no result for single word continuing loop")
		}

		if len(singleSearchResult) > 0 {
			for _, searchResult := range singleSearchResult {
				translation, _ := d.parseDictResults(searchResult)
				doNotAdd := false
				for _, res := range results.Results {
					if res.Translation == translation {
						doNotAdd = true
						break
					}
				}

				if doNotAdd {
					continue
				}

				result := models.Result{
					Word:        word,
					Rule:        "preposition",
					RootWord:    searchResult.Greek,
					Translation: translation,
				}
				results.Results = append(results.Results, result)
			}
		}
	}

	if len(results.Results) > 1 {
		lastResult := models.Result{
			Word:        "",
			Rule:        "",
			RootWord:    "",
			Translation: "",
		}
		for i, result := range results.Results {
			if result.Rule != "preposition" {
				if lastResult.Rule == result.Rule && result.Translation == lastResult.Translation {
					results.RemoveIndex(i)
				}
			}
			if result.Translation == "" {
				results.RemoveIndex(i)
			}
			lastResult = result
		}
	}
	return &results, nil
}

func (d *DionysosHandler) searchForDeclensions(word string) (*models.FoundRules, error) {
	var foundRules models.FoundRules

	for _, declension := range d.Config.DeclensionConfig.Declensions {
		var contract bool
		switch declension.Type {
		case "past":
			contract = true
		case "irregular":
			rules := d.loopOverIrregularVerbs(word, declension.Declensions)
			for _, rule := range rules.Rules {
				foundRules.Rules = append(foundRules.Rules, rule)
			}
			continue
		default:
			contract = false
		}

		for _, declensionForm := range declension.Declensions {
			result := d.loopOverDeclensions(word, declensionForm, contract)
			if len(result.Rules) >= 1 {
				for _, rule := range result.Rules {
					inArray := seeIfStringIsInArray(rule.Rule, foundRules.Rules)
					if inArray {
						continue
					}
					foundRules.Rules = append(foundRules.Rules, rule)
				}
			}
		}
	}

	return &foundRules, nil
}

func (d *DionysosHandler) loopOverDeclensions(word string, form models.DeclensionElement, contraction bool) models.FoundRules {
	var declensions models.FoundRules

	rootCutOff := 0
	if contraction {
		rootCutOff = 1
	}
	trimmedLetters := d.removeAccents(strings.Replace(form.Declension, "-", "", -1))
	lengthOfDeclension := utf8.RuneCountInString(trimmedLetters)
	wordInRune := []rune(word)
	if lengthOfDeclension > len(wordInRune) {
		return declensions
	}

	lettersOfWord := d.removeAccents(string(wordInRune[len(wordInRune)-lengthOfDeclension:]))
	if lettersOfWord == trimmedLetters {
		rootOfWord := string(wordInRune[rootCutOff : len(wordInRune)-lengthOfDeclension])
		firstLetter := d.removeAccents(string(wordInRune[0]))
		var words []string
		for _, term := range form.SearchTerm {
			if contraction {
				legitimateStartLetters := []string{"η", "ε"}
				legitimate := false
				for _, startLetter := range legitimateStartLetters {
					if startLetter == firstLetter {
						legitimate = true
					}
				}

				if !legitimate {
					continue
				}
				if firstLetter == "η" {
					vowels := []string{"α", "ε"}
					for _, vowel := range vowels {
						searchTerm := fmt.Sprintf("%s%s%s", vowel, rootOfWord, term)
						words = append(words, searchTerm)
					}
					continue

				}
			}

			searchTerm := fmt.Sprintf("%s%s", rootOfWord, term)
			words = append(words, searchTerm)
		}

		declension := models.Rule{
			Rule:        form.RuleName,
			SearchTerms: words,
		}

		declensions.Rules = append(declensions.Rules, declension)
	}

	return declensions
}

func (d *DionysosHandler) loopOverIrregularVerbs(word string, declensions []models.DeclensionElement) models.FoundRules {
	var rules models.FoundRules
	strippedWord := d.removeAccents(word)

	for _, outcome := range declensions {
		strippedOutcomeWord := d.removeAccents(outcome.Declension)
		if strippedWord == strippedOutcomeWord {
			declension := models.Rule{
				Rule:        outcome.RuleName,
				SearchTerms: outcome.SearchTerm,
			}
			rules.Rules = append(rules.Rules, declension)
		}
	}

	return rules
}

func seeIfStringIsInArray(s string, slice []models.Rule) bool {
	for _, field := range slice {
		if field.Rule == s {
			return true
		}
	}

	return false
}
