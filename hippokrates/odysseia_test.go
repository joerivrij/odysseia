package hippokrates

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/kpango/glg"
	"github.com/odysseia/hippokrates/client"
	"github.com/odysseia/hippokrates/client/models"
	"os"
	"testing"
)

const (
	sokratesApi   = "sokrates"
	herodotosApi  = "herodotos"
	alexandrosApi = "alexandros"
	dionysiosApi  = "dionysios"
	ResponseBody  = "responseBody"
	ErrorBody     = "errorBody"
	ContextAuthor = "contextAuthor"
	AnswerBody    = "answerBody"
)

var baseConfig *client.ClientConfig

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func (l *odysseiaFixture) theIsRunning(service string) error {
	var response *models.Health
	var err error

	switch service {
	case alexandrosApi:
		response, err = l.clients.Alexandros().Health()
		if err != nil {
			return err
		}
	case sokratesApi:
		response, err = l.clients.Sokrates().Health()
		if err != nil {
			return err
		}
	case herodotosApi:
		response, err = l.clients.Herodotos().Health()
		if err != nil {
			return err
		}
	case dionysiosApi:
		response, err = l.clients.Dionysios().Health()
		if err != nil {
			return err
		}
	default:
	}

	if !response.Healthy {
		return fmt.Errorf("service was %v were a healthy status was expected", response.Healthy)
	}

	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {

		//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HIPPOKRATES
		glg.Info("\n __ __  ____  ____  ____   ___   __  _  ____    ____  ______    ___  _____\n|  |  ||    ||    \\|    \\ /   \\ |  |/ ]|    \\  /    ||      |  /  _]/ ___/\n|  |  | |  | |  o  )  o  )     ||  ' / |  D  )|  o  ||      | /  [_(   \\_ \n|  _  | |  | |   _/|   _/|  O  ||    \\ |    / |     ||_|  |_||    _]\\__  |\n|  |  | |  | |  |  |  |  |     ||     ||    \\ |  _  |  |  |  |   [_ /  \\ |\n|  |  | |  | |  |  |  |  |     ||  .  ||  .  \\|  |  |  |  |  |     |\\    |\n|__|__||____||__|  |__|   \\___/ |__|\\_||__|\\_||__|__|  |__|  |_____| \\___|\n                                                                          \n")
		glg.Info("\"ὄμνυμι Ἀπόλλωνα ἰητρὸν καὶ Ἀσκληπιὸν καὶ Ὑγείαν καὶ Πανάκειαν καὶ θεοὺς πάντας τε καὶ πάσας, ἵστορας ποιεύμενος, ἐπιτελέα ποιήσειν κατὰ δύναμιν καὶ κρίσιν ἐμὴν ὅρκον τόνδε καὶ συγγραφὴν τήνδε:\"")
		glg.Info("\"I swear by Apollo Healer, by Asclepius, by Hygieia, by Panacea, and by all the gods and goddesses, making them my witnesses, that I will carry out, according to my ability and judgment, this oath and this indenture.\"")
		glg.Info("starting test suite setup.....")
		glg.Debug("getting env variables and creating config")

		config, err := GetEnv()
		if err != nil {
			glg.Fatal("could not parse config")
		}

		baseConfig = config
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
	})

	odysseia, err := New(baseConfig)
	if err != nil {
		os.Exit(1)
	}

	//general
	ctx.Step(`^the "([^"]*)" is running$`, odysseia.theIsRunning)

	//alexandros
	ctx.Step(`^the word "([^"]*)" is queried$`, odysseia.theWordIsQueried)
	ctx.Step(`^the word "([^"]*)" is stripped of accents$`, odysseia.theWordIsStrippedOfAccents)
	ctx.Step(`^the partial "([^"]*)" is queried$`, odysseia.thePartialIsQueried)
	ctx.Step(`^the word "([^"]*)" is queried with an error$`, odysseia.theWordIsQueriedWithAnError)
	ctx.Step(`^the word "([^"]*)" should be included in the response$`, odysseia.theWordShouldBeIncludedInTheResponse)
	ctx.Step(`^the number of results should not exceed "([^"]*)"$`, odysseia.theNumberOfResultsShouldNotExceed)
	ctx.Step(`^an error containing "([^"]*)" is returned$`, odysseia.anErrorContainingIsReturned)

	//herodotos
	ctx.Step(`^a query is made for all authors$`, odysseia.aQueryIsMadeForAllAuthors)
	ctx.Step(`^the author "([^"]*)" should be included$`, odysseia.theAuthorShouldBeIncluded)
	ctx.Step(`^a query is made for all books by author "([^"]*)"$`, odysseia.aQueryIsMadeForAllBooksByAuthor)
	ctx.Step(`^a translation is returned$`, odysseia.aTranslationIsReturned)
	ctx.Step(`^an author and book combination is queried$`, odysseia.anAuthorAndBookCombinationIsQueried)
	ctx.Step(`^the number of authors should exceed "([^"]*)"$`, odysseia.theNumberOfAuthorsShouldExceed)
	ctx.Step(`^the book "([^"]*)" should be included$`, odysseia.theBookShouldBeIncluded)
	ctx.Step(`^the sentenceId should be longer than "([^"]*)"$`, odysseia.theSentenceIdShouldBeLongerThan)
	ctx.Step(`^the sentence should include non-ASCII \(Greek\) characters$`, odysseia.theSentenceShouldIncludeNonASCIIGreekCharacters)
	ctx.Step(`^a correctness percentage$`, odysseia.aCorrectnessPercentage)
	ctx.Step(`^a sentence with a translation$`, odysseia.aSentenceWithATranslation)

	//sokrates
	ctx.Step(`^a query is made for all methods$`, odysseia.aQueryIsMadeForAllMethods)
	ctx.Step(`^the method "([^"]*)" should be included$`, odysseia.theMethodShouldBeIncluded)
	ctx.Step(`^a random method is queried for categories$`, odysseia.aRandomMethodIsQueriedForCategories)
	ctx.Step(`^the number of methods should exceed "([^"]*)"$`, odysseia.theNumberOfMethodsShouldExceed)
	ctx.Step(`^a category should be returned$`, odysseia.aCategoryShouldBeReturned)
	ctx.Step(`^a random category is queried for the last chapter$`, odysseia.aRandomCategoryIsQueriedForTheLastChapter)
	ctx.Step(`^that chapter should be a number above (\d+)$`, odysseia.thatChapterShouldBeANumberAbove)
	ctx.Step(`^a new quiz question is requested$`, odysseia.aNewQuizQuestionIsRequested)
	ctx.Step(`^that question is answered with a "([^"]*)" answer$`, odysseia.thatQuestionIsAnsweredWithAAnswer)
	ctx.Step(`^the result should be "([^"]*)"$`, odysseia.theResultShouldBe)

	//dionysios
	ctx.Step(`^the grammar is checked for word "([^"]*)"$`, odysseia.theGrammarIsCheckedForWord)
	ctx.Step(`^the grammar for word "([^"]*)" is queried with an error$`, odysseia.theGrammarForWordIsQueriedWithAnError)
	ctx.Step(`^the declension "([^"]*)" should be included in the response$`, odysseia.theDeclensionShouldBeIncludedInTheResponse)
	ctx.Step(`^the number of results should be equal to or exceed "([^"]*)"$`, odysseia.theNumberOfResultsShouldBeEqualToOrExceed)
	ctx.Step(`^the number of translations should be equal to er exceed "([^"]*)"$`, odysseia.theNumberOfTranslationsShouldBeEqualToErExceed)
	ctx.Step(`^the number of declensions should be equal to or exceed "([^"]*)"$`, odysseia.theNumberOfDeclensionsShouldBeEqualToOrExceed)

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
		Paths:  []string{"features"},
	}

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}
