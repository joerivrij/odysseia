package hippokrates

import (
	"context"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"strings"
)

func (l *odysseiaFixture) theWordIsQueried(word string) error {
	response, err := l.clients.Alexandros().QueryWord(word)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ResponseBody, response)

	return nil
}

func (l *odysseiaFixture) theWordIsQueriedWithAnError(word string) error {
	_, err := l.clients.Alexandros().QueryWord(word)
	if err != nil {
		l.ctx = context.WithValue(l.ctx, ErrorBody, err.Error())
	}

	return nil
}

func (l *odysseiaFixture) thePartialIsQueried(partial string) error {
	response, err := l.clients.Alexandros().QueryWord(partial)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ResponseBody, response)

	return nil
}

func (l *odysseiaFixture) theWordIsStrippedOfAccents(word string) error {
	strippedWord := RemoveAccents(word)

	response, err := l.clients.Alexandros().QueryWord(strippedWord)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ResponseBody, response)

	return nil
}

func (l *odysseiaFixture) theWordShouldBeIncludedInTheResponse(searchTerm string) error {
	words := l.ctx.Value(ResponseBody).([]models.Meros)

	found := false

	for _, word := range words {
		if strings.Contains(word.Greek, searchTerm) {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("could not find searchterm %v in slice", searchTerm)
	}
	return nil
}

func (l *odysseiaFixture) theNumberOfResultsShouldNotExceed(results int) error {
	words := l.ctx.Value(ResponseBody).([]models.Meros)
	lengthOfResult := len(words)

	if lengthOfResult > results {
		return fmt.Errorf("number of results is %v were a max of %v was expected", lengthOfResult, results)
	}

	return nil
}

func (l *odysseiaFixture) anErrorContainingIsReturned(message string) error {
	errorText := l.ctx.Value(ErrorBody).(string)
	if !strings.Contains(errorText, message) {
		return fmt.Errorf("expected %v to contain %v", errorText, message)
	}

	return nil
}
