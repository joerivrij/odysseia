package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/herodotos/app"
	"github.com/odysseia/plato/elastic"
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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HERODOTOS
	glg.Info("\n __ __    ___  ____   ___   ___     ___   ______   ___   _____\n|  |  |  /  _]|    \\ /   \\ |   \\   /   \\ |      | /   \\ / ___/\n|  |  | /  [_ |  D  )     ||    \\ |     ||      ||     (   \\_ \n|  _  ||    _]|    /|  O  ||  D  ||  O  ||_|  |_||  O  |\\__  |\n|  |  ||   [_ |    \\|     ||     ||     |  |  |  |     |/  \\ |\n|  |  ||     ||  .  \\     ||     ||     |  |  |  |     |\\    |\n|__|__||_____||__|\\_|\\___/ |_____| \\___/   |__|   \\___/  \\___|\n                                                              \n")
	glg.Info("\"Ἡροδότου Ἁλικαρνησσέος ἱστορίης ἀπόδεξις ἥδε\"")
	glg.Info("\"This is the display of the inquiry of Herodotos of Halikarnassos\"")
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
