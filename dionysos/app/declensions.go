package app

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/models"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

type DeclensionHandler struct {
	BaseUrl            string
	Version            string
	ApiName            string
	SearchWordEndPoint string
	Index string
	ElasticClient elasticsearch.Client
}

func (d *DeclensionHandler) queryAlexandrosForPossibleMeaning(word string) ([]models.Meros, error) {
	var results []models.Meros

	u, err := url.Parse(d.BaseUrl)
	if err != nil {
		return nil, err
	}

	strippedWord := d.removeAccents(word)

	u.Path = path.Join(u.Path, d.ApiName, d.Version, d.SearchWordEndPoint)
	q := u.Query()
	q.Set("word", strippedWord)
	u.RawQuery = q.Encode()

	response, err := helpers.GetRequest(*u)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&results)
	glg.Debug(results)

	return results, nil
}

func (d *DeclensionHandler)removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func (d *DeclensionHandler) StartFindingRules(word string) (*DeclensionTranslationResults, error) {
	var results DeclensionTranslationResults

	declensions, err := d.searchForDeclensions(word)
	if err != nil {
		return nil, err
	}

	if len(declensions.Rules) > 0 {
		for _, declension := range declensions.Rules {
			if len(declension.SearchTerms) > 0 {
				for _, term := range declension.SearchTerms {
					alexandrosHits, err := d.queryAlexandrosForPossibleMeaning(term)
					if err != nil {
						return nil, err
					}
					var translation string
					var article string

					if len(alexandrosHits) > 0 {
						translation = alexandrosHits[0].English
						greek := strings.Split(alexandrosHits[0].Greek, ",")
						if len(greek) > 1 {
							article = strings.Replace(greek[1], " ", "", -1)
						}
					}

					result := Result{
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

	if len(results.Results) > 1 {
		for i, result := range results.Results {
			if result.Translation == "" {
				results.RemoveIndex(i)
			}
		}
	}
	return &results, nil
}

func (d *DeclensionHandler) searchForDeclensions(word string) (*models.FoundRules, error) {
	var local bool

	localString := os.Getenv("LOCAL")
	if localString == "" {
		local = true
	} else {
		local = false
	}

	var declensions models.Declensions
	if local {
		var basePath string
		cw, _ := os.Getwd()
		if !strings.Contains(cw, "dionysos") {
			basePath = fmt.Sprintf("%s/%s/%s", cw, "dionysos", "app")
		} else {
			basePath = cw
		}
		filePath := fmt.Sprintf("%s/rules.json", basePath)
		jsonFile, err := os.Open(filePath)

		if err != nil {
			return nil, err
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)
		declensions, err = models.UnmarshalDeclensions(byteValue)
		if err != nil {
			return nil, err
		}
	} else {
		response, err := elastic.QueryWithMatchAll(d.ElasticClient, d.Index)

		if err != nil {
			return nil, err
		}

		jsonHit, _ := json.Marshal(response.Hits.Hits[0].Source)
		declensions, _ = models.UnmarshalDeclensions(jsonHit)

	}

	var foundRules models.FoundRules

	firstDeclensionForms := d.loopOverDeclensions(word, declensions.Declensions.FirstDeclension)
	for _, form := range firstDeclensionForms.Rules {
		rule := models.Rule{
			Form:        "noun",
			Declension:  "first",
			Rule:        form.Rule,
			SearchTerms: form.SearchTerms,
		}
		foundRules.Rules = append(foundRules.Rules, rule)
	}
	secondDeclensionForms := d.loopOverDeclensions(word, declensions.Declensions.SecondDeclension)
	for _, form := range secondDeclensionForms.Rules {
		rule := models.Rule{
			Form:        "noun",
			Declension:  "second",
			Rule:        form.Rule,
			SearchTerms: form.SearchTerms,
		}
		foundRules.Rules = append(foundRules.Rules, rule)
	}

	return &foundRules, nil
}

func (d *DeclensionHandler) loopOverDeclensions(word string, form []models.DeclensionElement) models.FoundRules {
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
