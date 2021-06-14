package hippokrates

import (
	"context"
)

func (l *LexikoFixture)aNewSentenceIsRequestedForAuthor(author string) error {
	response, err := l.herodotos.CreateSentence(author)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, StatusCode, response.StatusCode)

	return nil
}
