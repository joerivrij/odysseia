package app

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
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
	response, err := elastic.QueryWithMatch(d.Config.ElasticClient, d.Config.SecondaryIndex, term, strippedWord)

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
						//this should be handled in a new service
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

	if len(results.Results) > 1 {
		lastRule := ""
		for i, result := range results.Results {
			if result.Rule != "preposition" {
				if lastRule == result.Rule {
					results.RemoveIndex(i)
				}
			}
			if result.Translation == "" {
				results.RemoveIndex(i)
			}
			lastRule = result.Rule
		}
	}
	return &results, nil
}

func (d *DionysosHandler) searchForDeclensions(word string) (*models.FoundRules, error) {
	var foundRules models.FoundRules

	firstDeclensionForms := d.loopOverDeclensions(word, d.Config.DeclensionConfig.FirstDeclension.Declensions)
	for _, form := range firstDeclensionForms.Rules {
		rule := models.Rule{
			Form:        d.Config.DeclensionConfig.FirstDeclension.Type,
			Declension:  d.Config.DeclensionConfig.FirstDeclension.Name,
			Rule:        form.Rule,
			SearchTerms: form.SearchTerms,
		}
		foundRules.Rules = append(foundRules.Rules, rule)
	}
	secondDeclensionForms := d.loopOverDeclensions(word, d.Config.DeclensionConfig.SecondDeclension.Declensions)
	for _, form := range secondDeclensionForms.Rules {
		rule := models.Rule{
			Form:        d.Config.DeclensionConfig.SecondDeclension.Type,
			Declension:  d.Config.DeclensionConfig.SecondDeclension.Name,
			Rule:        form.Rule,
			SearchTerms: form.SearchTerms,
		}
		foundRules.Rules = append(foundRules.Rules, rule)
	}

	return &foundRules, nil
}

func (d *DionysosHandler) loopOverDeclensions(word string, form []models.DeclensionElement) models.FoundRules {
	var declensions models.FoundRules

	for _, outcome := range form {
		trimmedLetters := strings.Replace(outcome.Declension, "-", "", -1)
		lengthOfDeclension := utf8.RuneCountInString(trimmedLetters)
		wordInRune := []rune(word)
		lettersOfWord := string(wordInRune[len(wordInRune)-lengthOfDeclension:])
		if lettersOfWord == trimmedLetters {
			rootOfWord := string(wordInRune[0 : len(wordInRune)-lengthOfDeclension])
			var words []string
			for _, term := range outcome.SearchTerm {
				searchTerm := fmt.Sprintf("%s%s", rootOfWord, term)
				words = append(words, searchTerm)
			}

			if len(declensions.Rules) > 0 {
				inArray := seeIfStringIsInArray(outcome.RuleName, declensions.Rules)

				if inArray {
					continue
				}
				declension := models.Rule{
					Rule:        outcome.RuleName,
					SearchTerms: words,
				}
				declensions.Rules = append(declensions.Rules, declension)
			} else {
				declension := models.Rule{
					Rule:        outcome.RuleName,
					SearchTerms: words,
				}
				declensions.Rules = append(declensions.Rules, declension)
			}
		}
	}

	return declensions
}

func seeIfStringIsInArray(s string, slice []models.Rule) bool {
	for _, field := range slice {
		if field.Rule == s {
			return true
		}
	}

	return false
}
