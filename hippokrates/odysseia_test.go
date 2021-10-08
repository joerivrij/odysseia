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
	dionysosApi = "dionysos"
	StatusCode = "statusCode"
	CreateQuestionStatusCode = "createQuestionStatusCode"
	CreateSentenceStatusCode = "createSentenceStatusCode"
	QueryWordStatusCode = "queryWordStatusCode"
)

var ALEXANDROS_URL string
var HERODOTOS_URL string
var SOKRATES_URL string
var DIONYSOS_URL string

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func (l *odysseiaFixture)theIsRunning(service string) error {
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
	case dionysosApi:
		response, err = l.dionysos.Health()
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

func (l *odysseiaFixture)theResponseCodeShouldBe(code int) error {
	statusCode := l.ctx.Value(StatusCode).(int)
	if statusCode != code {
		return fmt.Errorf("code was %d where %d was expected", statusCode, code)
	}

	return nil
}


func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		alexandrosUrl := envflag.String("ALEXANDROS_URL", "http://minikube-odysseia.test", "alexandros base url")
		herodotosUrl := envflag.String("HERODOTOS_URL", "http://minikube-odysseia.test", "herodotos base url")
		sokratesUrl := envflag.String("SOKRATES_URL", "http://minikube-odysseia.test", "sokrates base url")
		dionysosUrl := envflag.String("DIONYSOS_URL", "http://minikube-odysseia.test", "dionysos base url")

		envflag.Parse()
		flag.Parse()

		ALEXANDROS_URL = *alexandrosUrl
		HERODOTOS_URL = *herodotosUrl
		SOKRATES_URL = *sokratesUrl
		DIONYSOS_URL = *dionysosUrl
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
	})

	odysseia, err := New(ALEXANDROS_URL, HERODOTOS_URL, SOKRATES_URL, DIONYSOS_URL, sokratesApi, herodotosApi, alexandrosApi, dionysosApi)
	if err != nil {
		os.Exit(1)
	}

	//general
	ctx.Step(`^the "([^"]*)" is running$`, odysseia.theIsRunning)
	ctx.Step(`^the responseCode should be "([^"]*)"$`, odysseia.theResponseCodeShouldBe)

	//alexandros
	ctx.Step(`^the word "([^"]*)" is queried$`, odysseia.theWordIsQueried)

	//herodotos
	ctx.Step(`^a new sentence is requested for author "([^"]*)"$`, odysseia.aNewSentenceIsRequestedForAuthor)

	//sokrates
	ctx.Step(`^a new question is requested with category "([^"]*)" and chapter "([^"]*)"$`, odysseia.aNewQuestionIsRequestedWithCategoryAndChapter)

	//dionysos
	ctx.Step(`^the grammar is checked for word "([^"]*)"$`, odysseia.theGrammarIsCheckedForWord)
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
