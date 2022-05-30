package hippokrates

import (
	"context"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
)

func (l *odysseiaFixture) theGrammarForWordIsQueriedWithAnError(word string) error {
	_, err := l.clients.Dionysios().CheckGrammar(word)
	if err != nil {
		l.ctx = context.WithValue(l.ctx, ErrorBody, err.Error())
	}

	return nil
}

func (l *odysseiaFixture) theGrammarIsCheckedForWord(word string) error {
	response, err := l.clients.Dionysios().CheckGrammar(word)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ResponseBody, response)

	return nil
}

func (l *odysseiaFixture) theDeclensionShouldBeIncludedInTheResponse(declension string) error {
	declensions := l.ctx.Value(ResponseBody).(*models.DeclensionTranslationResults)

	found := false
	for _, decResult := range declensions.Results {
		if decResult.Rule == declension {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("could not find declension %v in slice", declension)
	}
	return nil
}

func (l *odysseiaFixture) theNumberOfResultsShouldBeEqualToOrExceed(numberOfResults int) error {
	declensions := l.ctx.Value(ResponseBody).(*models.DeclensionTranslationResults)
	lengthOfResults := len(declensions.Results)
	if lengthOfResults < numberOfResults {
		return fmt.Errorf("expected results to be equal to or more than %v but was %v", numberOfResults, lengthOfResults)
	}

	return nil
}

func (l *odysseiaFixture) theNumberOfTranslationsShouldBeEqualToErExceed(numberOfTranslations int) error {
	declensions := l.ctx.Value(ResponseBody).(*models.DeclensionTranslationResults)
	var translations []string
	for _, result := range declensions.Results {
		inTranslation := false
		for _, translation := range translations {
			if translation == result.Translation {
				inTranslation = true
			}
		}

		if !inTranslation {
			translations = append(translations, result.Translation)
		}
	}

	lengthOTranslations := len(translations)
	if lengthOTranslations < numberOfTranslations {
		return fmt.Errorf("expected translation results to be equal to or more than %v but was %v", numberOfTranslations, lengthOTranslations)
	}

	return nil
}

func (l *odysseiaFixture) theNumberOfDeclensionsShouldBeEqualToOrExceed(numberOfDeclensions int) error {
	declensions := l.ctx.Value(ResponseBody).(*models.DeclensionTranslationResults)

	var declensionRules []string
	for _, result := range declensions.Results {
		inTranslation := false
		for _, translation := range declensionRules {
			if translation == result.Rule {
				inTranslation = true
			}
		}

		if !inTranslation {
			declensionRules = append(declensionRules, result.Rule)
		}
	}

	lengthOfDeclensions := len(declensionRules)
	if lengthOfDeclensions < numberOfDeclensions {
		return fmt.Errorf("expected declension results to be equal to or more than %v but was %v", numberOfDeclensions, lengthOfDeclensions)
	}

	return nil
}
