package hippokrates

import "context"

func (l *odysseiaFixture)theGrammarIsCheckedForWord(word string) error {
	response, err := l.dionysos.CheckGrammar(word)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, StatusCode, response.StatusCode)

	return nil
}
