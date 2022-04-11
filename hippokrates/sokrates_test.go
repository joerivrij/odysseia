package hippokrates

import (
	"context"
	"fmt"
	"github.com/odysseia/hippokrates/client/models"
	"strconv"
)

const (
	ContextCategories string = "contextCategories"
	ContextChapter    string = "contextChapter"
	ContextMethod     string = "contextMethod"
	ContextAnswer     string = "contextAnswer"
)

func (l *odysseiaFixture) aQueryIsMadeForAllMethods() error {
	response, err := l.clients.Sokrates().Methods()
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ResponseBody, response)

	return nil
}

func (l *odysseiaFixture) aRandomMethodIsQueriedForCategories() error {
	methods := l.ctx.Value(ResponseBody).(*models.Methods)

	randNumber := GenerateRandomNumber(len(methods.Method))
	method := methods.Method[randNumber].Method

	categories, err := l.clients.Sokrates().Categories(method)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ContextCategories, categories)
	l.ctx = context.WithValue(l.ctx, ContextMethod, method)

	return nil
}

func (l *odysseiaFixture) aRandomCategoryIsQueriedForTheLastChapter() error {
	categories := l.ctx.Value(ContextCategories).(*models.Categories)
	method := l.ctx.Value(ContextMethod).(string)

	randNumber := GenerateRandomNumber(len(categories.Category))
	category := categories.Category[randNumber].Category

	lastChapter, err := l.clients.Sokrates().LastChapter(method, category)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ContextChapter, lastChapter)

	return nil
}

func (l *odysseiaFixture) aNewQuizQuestionIsRequested() error {
	methods, err := l.clients.Sokrates().Methods()
	if err != nil {
		return err
	}

	randomMethod := GenerateRandomNumber(len(methods.Method))
	method := methods.Method[randomMethod].Method

	categories, err := l.clients.Sokrates().Categories(method)
	if err != nil {
		return err
	}

	randomCategory := GenerateRandomNumber(len(categories.Category))
	category := categories.Category[randomCategory].Category

	lastChapter, err := l.clients.Sokrates().LastChapter(method, category)
	if err != nil {
		return err
	}

	randomChapter := GenerateRandomNumber(int(lastChapter.LastChapter)) + 1

	chapter := strconv.Itoa(randomChapter)
	quizQuestion, err := l.clients.Sokrates().CreateQuestion(method, category, chapter)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ContextCategories, category)
	l.ctx = context.WithValue(l.ctx, ResponseBody, quizQuestion)

	return nil
}

func (l *odysseiaFixture) thatQuestionIsAnsweredWithAAnswer(correctAnswer string) error {
	quiz := l.ctx.Value(ResponseBody).(models.QuizResponse)
	category := l.ctx.Value(ContextCategories).(string)

	request := models.CheckAnswerRequest{
		QuizWord:       quiz[0],
		AnswerProvided: "",
		Category:       category,
	}

	parsedAnswer, err := strconv.ParseBool(correctAnswer)
	if err != nil {
		return err
	}

	if parsedAnswer {
		request.AnswerProvided = quiz[1]
	} else {
		request.AnswerProvided = quiz[2]
	}

	answer, err := l.clients.Sokrates().Answer(request)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ContextAnswer, answer.Correct)

	return nil
}

func (l *odysseiaFixture) theMethodShouldBeIncluded(method string) error {
	methods := l.ctx.Value(ResponseBody).(*models.Methods)

	found := false

	for _, result := range methods.Method {
		if result.Method == method {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("could not find book %v in slice", method)
	}

	return nil
}

func (l *odysseiaFixture) theNumberOfMethodsShouldExceed(results int) error {
	methods := l.ctx.Value(ResponseBody).(*models.Methods)
	numberOfMethods := len(methods.Method)
	if numberOfMethods <= results {
		return fmt.Errorf("expected results to be equal to or more than %v but was %v", results, numberOfMethods)
	}

	return nil
}

func (l *odysseiaFixture) aCategoryShouldBeReturned() error {
	categories := l.ctx.Value(ContextCategories).(*models.Categories)

	if len(categories.Category) == 0 {
		return fmt.Errorf("expected categories to be returned but non were found")
	}

	return nil
}

func (l *odysseiaFixture) thatChapterShouldBeANumberAbove(number int) error {
	lastChapter := l.ctx.Value(ContextChapter).(*models.LastChapterResponse)
	if lastChapter.LastChapter < int64(number) {
		return fmt.Errorf("expected lastchapter to be higher than %v but was %v", number, lastChapter.LastChapter)
	}

	return nil
}

func (l *odysseiaFixture) theResultShouldBe(correct string) error {
	answer := l.ctx.Value(ContextAnswer).(bool)

	parsedCorrectness, err := strconv.ParseBool(correct)
	if err != nil {
		return err
	}

	if answer != parsedCorrectness {
		return fmt.Errorf("expected answer %v to be equal to correctness %v", answer, correct)
	}
	return nil
}
