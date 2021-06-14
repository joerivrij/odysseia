package hippokrates

import (
	"context"
)

func (l *LexikoFixture)theWordIsQueried(word string) error {
	response, err := l.alexandros.QueryWord(word)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, StatusCode, response.StatusCode)

	return nil
}
