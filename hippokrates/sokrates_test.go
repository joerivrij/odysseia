package hippokrates

import (
	"context"
)

func (l *LexikoFixture)aNewQuestionIsRequestedWithCategoryAndChapter(category, chapter string) error {
	response, err := l.sokrates.CreateQuestion(category, chapter)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, StatusCode, response.StatusCode)

	return nil
}
