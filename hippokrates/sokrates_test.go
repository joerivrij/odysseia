package hippokrates

import (
	"context"
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/ianschenck/envflag"
	"net/http"
	"os"
	"testing"
)

const (
	sokratesApi = "sokrates"
	CreateQuestionStatusCode = "createQuestionStatusCode"
)

var SOKRATES_URL string

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func (s *SokratesFixture)theIsRunning(service string) error {
	var response *http.Response
	var err error
	expectedCode := 200

	switch service {
	case "sokrates":
		response, err = s.sokrates.Health()
		if err != nil {
			return err
		}
	default:
	}

	if response.StatusCode != expectedCode {
		return fmt.Errorf("code was %d where %d was expected", response.StatusCode, expectedCode)
	}

	return nil
}

func (s *SokratesFixture)aNewQuestionIsRequestedWithCategoryAndChapter(category, chapter string) error {
	response, err := s.sokrates.CreateQuestion(category, chapter)
	if err != nil {
		return err
	}

	s.ctx = context.WithValue(s.ctx, CreateQuestionStatusCode, response.StatusCode)

	return nil
}


func (s *SokratesFixture)theResponseCodeShouldBe(code int) error {
	statusCode := s.ctx.Value(CreateQuestionStatusCode).(int)
	if statusCode != code {
		return fmt.Errorf("code was %d where %d was expected", statusCode, code)
	}

	return nil
}


func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		sokratesUrl := envflag.String("SOKRATES_URL", "http://minikube-lexiko.test", "sokrates base url")

		envflag.Parse()
		flag.Parse()

		SOKRATES_URL = *sokratesUrl
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
	})

	sokrates, err := New(SOKRATES_URL, sokratesApi)
	if err != nil {
		os.Exit(1)
	}

	ctx.Step(`^a new question is requested with category "([^"]*)" and chapter "([^"]*)"$`, sokrates.aNewQuestionIsRequestedWithCategoryAndChapter)
	ctx.Step(`^the "([^"]*)" is running$`, sokrates.theIsRunning)
	ctx.Step(`^the responseCode should be "([^"]*)"$`, sokrates.theResponseCodeShouldBe)
}

func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opts := godog.Options{
		Format: format,
		Paths:     []string{"features"},
	}

	status := godog.TestSuite{
		Name: "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	os.Exit(status)
}
