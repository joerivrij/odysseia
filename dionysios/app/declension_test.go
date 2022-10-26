package app

import (
	"github.com/odysseia-greek/plato/models"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckGrammarEndPointIrregularVerb(t *testing.T) {
	numberOfRules := 1

	t.Run("HappyPathIrregularVerb", func(t *testing.T) {
		searchWord := "ᾖσαν"
		expected := "3th plural - impf - ind - act"
		expectedSearchResult := "εἰμί"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}

		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case "irregular":
				rules := handler.loopOverIrregularVerbs(searchWord, declension.Declensions)
				for _, rule := range rules.Rules {
					foundRules.Rules = append(foundRules.Rules, rule)
				}

			default:
				continue
			}
		}

		assert.Nil(t, err)
		assert.True(t, len(foundRules.Rules) == numberOfRules)
		expectedRuleFound := false
		for _, rule := range foundRules.Rules {
			if rule.Rule == expected {
				expectedRuleFound = true
			}
			assert.Equal(t, expectedSearchResult, rule.SearchTerms[0])
		}
		assert.True(t, expectedRuleFound)
	})

	t.Run("FullPathIrregularVerb", func(t *testing.T) {
		searchWord := "ᾖσαν"
		expected := "3th plural - impf - ind - act"
		expectedSearchResult := "εἰμί"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{
			Config: &configs.DionysiosConfig{DeclensionConfig: *declensionConfig},
		}

		foundRules, err := handler.searchForDeclensions(searchWord)
		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.True(t, len(foundRules.Rules) == 4)
		expectedRuleFound := false
		for _, rule := range foundRules.Rules {
			if rule.Rule == expected {
				expectedRuleFound = true
				assert.Equal(t, expectedSearchResult, rule.SearchTerms[0])
			}
		}
		assert.True(t, expectedRuleFound)
	})
}

func TestDeclensionImperfectumResult(t *testing.T) {
	numberOfRules := 2
	contraction := true

	t.Run("HappyPathImperfectum", func(t *testing.T) {
		searchWord := "ἔφερον"
		expected := "1st sing - impf - ind - act"
		expectedSearchResult := "φερω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}

		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case "imperfect":
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		expectedRuleFound := false
		for _, rule := range foundRules.Rules {
			if rule.Rule == expected {
				expectedRuleFound = true
			}
			assert.Equal(t, expectedSearchResult, rule.SearchTerms[0])
		}
		assert.True(t, expectedRuleFound)

	})
}

func TestDeclensionAoristResult(t *testing.T) {
	numberOfRules := 1
	multipleSearchResults := 4
	contraction := true
	name := "firstAorist"

	t.Run("HappyPathFirstAoristPsi", func(t *testing.T) {
		searchWord := "ἔγρᾰψᾰ"
		expected := "1st sing - aorist - ind - act"
		expectedSearchResult := "γρᾰφω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}
		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case name:
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		assert.Equal(t, expected, foundRules.Rules[0].Rule)
		assert.Equal(t, expectedSearchResult, foundRules.Rules[0].SearchTerms[0])
	})

	t.Run("HappyPathFirstAoristSigma", func(t *testing.T) {
		searchWord := "ἐλύσαμεν"
		expected := "1st plural - aorist - ind - act"
		expectedSearchResult := "λύω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}
		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case name:
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		assert.Equal(t, expected, foundRules.Rules[0].Rule)
		assert.Equal(t, expectedSearchResult, foundRules.Rules[0].SearchTerms[0])
	})

	t.Run("HappyPathFirstAoristKappa", func(t *testing.T) {
		searchWord := "ἔπλέξεν"
		expected := "3th sing - aorist - ind - act"
		expectedSearchResult := "πλέκω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}
		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case name:
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		assert.Equal(t, expected, foundRules.Rules[0].Rule)
		assert.Equal(t, multipleSearchResults, len(foundRules.Rules[0].SearchTerms))
		found := false
		for _, searchTerm := range foundRules.Rules[0].SearchTerms {
			if searchTerm == expectedSearchResult {
				found = true
			}
		}

		assert.True(t, found)
	})

	t.Run("HappyPathFirstAoristSigmaKappa", func(t *testing.T) {
		searchWord := "ἐδῐδᾰ́ξᾰτε"
		expected := "2nd plural - aorist - ind - act"
		expectedSearchResult := "δῐδᾰ́σκω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}
		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case name:
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}
		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		assert.Equal(t, expected, foundRules.Rules[0].Rule)
		assert.Equal(t, multipleSearchResults, len(foundRules.Rules[0].SearchTerms))
		found := false
		for _, searchTerm := range foundRules.Rules[0].SearchTerms {
			if searchTerm == expectedSearchResult {
				found = true
			}
		}

		assert.True(t, found)
	})

	t.Run("HappyPathFirstAoristGamma", func(t *testing.T) {
		searchWord := "ἔλεξᾰν"
		expected := "3th plural - aorist - ind - act"
		expectedSearchResult := "λεγω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}
		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case name:
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		assert.Equal(t, expected, foundRules.Rules[0].Rule)
		assert.Equal(t, multipleSearchResults, len(foundRules.Rules[0].SearchTerms))
		found := false
		for _, searchTerm := range foundRules.Rules[0].SearchTerms {
			if searchTerm == expectedSearchResult {
				found = true
			}
		}

		assert.True(t, found)
	})

	t.Run("HappyPathFirstAoristChiWithEta", func(t *testing.T) {
		searchWord := "ἦρξᾰς"
		expected := "2nd sing - aorist - ind - act"
		expectedSearchResult := "αρχω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}
		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Name {
			case name:
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		assert.Equal(t, expected, foundRules.Rules[0].Rule)
		assert.Equal(t, 2*multipleSearchResults, len(foundRules.Rules[0].SearchTerms))
		found := false
		for _, searchTerm := range foundRules.Rules[0].SearchTerms {
			if searchTerm == expectedSearchResult {
				found = true
			}
		}

		assert.True(t, found)
	})
}

func TestDeclensionParticiplesResult(t *testing.T) {
	numberOfRules := 1
	contraction := false

	t.Run("HappyPathParticpleMascSingNom", func(t *testing.T) {
		searchWord := "λυων"
		expected := "pres act part - sing - masc - nom"
		expectedSearchResult := "λυω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}

		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Type {
			case "participia":
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		expectedRuleFound := false
		for _, rule := range foundRules.Rules {
			if rule.Rule == expected {
				expectedRuleFound = true
			}
			assert.Equal(t, expectedSearchResult, rule.SearchTerms[0])
		}
		assert.True(t, expectedRuleFound)
	})

	t.Run("HappyPathParticpleFemDatPlural", func(t *testing.T) {
		searchWord := "λυοὐσαις"
		expected := "pres act part - plural - fem - dat"
		expectedSearchResult := "λυω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}

		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Type {
			case "participia":
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, numberOfRules, len(foundRules.Rules))
		expectedRuleFound := false
		for _, rule := range foundRules.Rules {
			if rule.Rule == expected {
				expectedRuleFound = true
			}
			assert.Equal(t, expectedSearchResult, rule.SearchTerms[0])
		}
		assert.True(t, expectedRuleFound)
	})

	t.Run("HappyPathParticpleNeutGenSing", func(t *testing.T) {
		searchWord := "λυὀντος"
		expected := "pres act part - sing - neut - gen"
		expectedSearchResult := "λυω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{}

		var foundRules models.FoundRules

		for _, declension := range declensionConfig.Declensions {
			switch declension.Type {
			case "participia":
				for _, element := range declension.Declensions {
					rules := handler.loopOverDeclensions(searchWord, element, contraction)
					for _, rule := range rules.Rules {
						foundRules.Rules = append(foundRules.Rules, rule)
					}
				}

			default:
				continue
			}
		}

		assert.Equal(t, 2, len(foundRules.Rules))
		expectedRuleFound := false
		for _, rule := range foundRules.Rules {
			if rule.Rule == expected {
				expectedRuleFound = true
			}
			assert.Equal(t, expectedSearchResult, rule.SearchTerms[0])
		}
		assert.True(t, expectedRuleFound)
	})

	t.Run("HappyPathParticpleAndVerbaResult", func(t *testing.T) {
		searchWord := "λυουσι"
		expected := "pres act part - plural - masc - dat"
		expectedVerba := "3th plural - pres - ind - act"
		expectedSearchResult := "λυω"

		declensionConfig, err := QueryRuleSet(nil, "dionysios")
		assert.Nil(t, err)

		handler := DionysosHandler{
			Config: &configs.DionysiosConfig{DeclensionConfig: *declensionConfig},
		}

		foundRules, err := handler.searchForDeclensions(searchWord)
		assert.Nil(t, err)

		assert.Equal(t, 3, len(foundRules.Rules))
		expectedRuleFound := false
		expectedVerbaFound := false
		for _, rule := range foundRules.Rules {
			if rule.Rule == expected {
				expectedRuleFound = true
			}

			if rule.Rule == expectedVerba {
				expectedVerbaFound = true
			}
			assert.Equal(t, expectedSearchResult, rule.SearchTerms[0])
		}
		assert.True(t, expectedRuleFound)
		assert.True(t, expectedVerbaFound)
	})
}
