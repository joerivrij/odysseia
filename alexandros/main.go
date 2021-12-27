package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/alexandros/app"
	"github.com/odysseia/plato/configuration"
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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ALEXANDROS
	glg.Info("\n  ____  _        ___  __ __   ____  ____   ___    ____   ___   _____\n /    || |      /  _]|  |  | /    ||    \\ |   \\  |    \\ /   \\ / ___/\n|  o  || |     /  [_ |  |  ||  o  ||  _  ||    \\ |  D  )     (   \\_ \n|     || |___ |    _]|_   _||     ||  |  ||  D  ||    /|  O  |\\__  |\n|  _  ||     ||   [_ |     ||  _  ||  |  ||     ||    \\|     |/  \\ |\n|  |  ||     ||     ||  |  ||  |  ||  |  ||     ||  .  \\     |\\    |\n|__|__||_____||_____||__|__||__|__||__|__||_____||__|\\_|\\___/  \\___|\n                                                                    \n")
	glg.Info("\"ὅτι τοῦ κρατεῖν πέρας ἡμῖν ἐστι τὸ μὴ ταὐτὰ ποιεῖν τοῖς κεκρατημένοις;’\"")
	glg.Info("\"Know ye not,’ said he, ‘that the end and object of conquest is to avoid doing the same thing as the conquered?\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")

	confManager, err := configuration.NewConfig()
	if err != nil {
		glg.Error(err)
		glg.Fatal("unable to fetch configuration")
	}
	config, err := app.GetFromConfig(confManager)
	if err != nil {
		glg.Error(err)
		glg.Fatal("unable to fetch configuration")
	}

	srv := app.InitRoutes(*config)

	glg.Infof("%s : %s", "running on port", port)
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
