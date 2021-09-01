package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/sokrates/app"
	"net/http"
	"os"
)

const standardPort = ":5000"

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=SOKRATES
	glg.Info("\n  _____  ___   __  _  ____    ____  ______    ___  _____\n / ___/ /   \\ |  |/ ]|    \\  /    ||      |  /  _]/ ___/\n(   \\_ |     ||  ' / |  D  )|  o  ||      | /  [_(   \\_ \n \\__  ||  O  ||    \\ |    / |     ||_|  |_||    _]\\__  |\n /  \\ ||     ||     ||    \\ |  _  |  |  |  |   [_ /  \\ |\n \\    ||     ||  .  ||  .  \\|  |  |  |  |  |     |\\    |\n  \\___| \\___/ |__|\\_||__|\\_||__|__|  |__|  |_____| \\___|\n                                                        \n")
	glg.Info("\"ἓν οἶδα ὅτι οὐδὲν οἶδα\"")
	glg.Info("\"I know one thing, that I know nothing\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")

	esClient, err := elastic.CreateElasticClientFromEnvVariables()
	if err != nil {
		glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
	}

	healthy, config := app.Get(200, esClient)
	if !healthy {
		glg.Fatal("death has found me")
	}

	srv := app.InitRoutes(*config)

	glg.Infof("%s : %s", "running on port", port)
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
