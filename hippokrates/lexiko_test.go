package hippokrates

import (
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
	herodotosApi = "herodotos"
	alexandrosApi = "alexandros"
	StatusCode = "statusCode"
	CreateQuestionStatusCode = "createQuestionStatusCode"
	CreateSentenceStatusCode = "createSentenceStatusCode"
	QueryWordStatusCode = "queryWordStatusCode"
)

var BASE_URL string

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func (l *LexikoFixture)theIsRunning(service string) error {
	var response *http.Response
	var err error
	expectedCode := 200

	switch service {
	case sokratesApi:
		response, err = l.sokrates.Health()
		if err != nil {
			return err
		}
	case herodotosApi:
		response, err = l.herodotos.Health()
		if err != nil {
			return err
		}
	case alexandrosApi:
		response, err = l.alexandros.Health()
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

func (l *LexikoFixture)theResponseCodeShouldBe(code int) error {
	statusCode := l.ctx.Value(StatusCode).(int)
	if statusCode != code {
		return fmt.Errorf("code was %d where %d was expected", statusCode, code)
	}

	return nil
}


func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		baseUrl := envflag.String("BASE_URL", "http://minikube-lexiko.test", "Lexiko base url")

		envflag.Parse()
		flag.Parse()

		BASE_URL = *baseUrl
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
	})

	lexiko, err := New(BASE_URL, sokratesApi, herodotosApi, alexandrosApi)
	if err != nil {
		os.Exit(1)
	}

	//general
	ctx.Step(`^the "([^"]*)" is running$`, lexiko.theIsRunning)
	ctx.Step(`^the responseCode should be "([^"]*)"$`, lexiko.theResponseCodeShouldBe)

	//alexandros
	ctx.Step(`^the word "([^"]*)" is queried$`, lexiko.theWordIsQueried)

	//herodotos
	ctx.Step(`^a new sentence is requested for author "([^"]*)"$`, lexiko.aNewSentenceIsRequestedForAuthor)

	//sokrates
	ctx.Step(`^a new question is requested with category "([^"]*)" and chapter "([^"]*)"$`, lexiko.aNewQuestionIsRequestedWithCategoryAndChapter)
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
