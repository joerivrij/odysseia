package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia/ptolemaios/app"
	"net/http"
	"os"
)

const standardPort = ":5001"

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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PTOLEMAIOS
	glg.Info("\n ____  ______   ___   _        ___  ___ ___   ____  ____  ___   _____\n|    \\|      | /   \\ | |      /  _]|   |   | /    ||    |/   \\ / ___/\n|  o  )      ||     || |     /  [_ | _   _ ||  o  | |  ||     (   \\_ \n|   _/|_|  |_||  O  || |___ |    _]|  \\_/  ||     | |  ||  O  |\\__  |\n|  |    |  |  |     ||     ||   [_ |   |   ||  _  | |  ||     |/  \\ |\n|  |    |  |  |     ||     ||     ||   |   ||  |  | |  ||     |\\    |\n|__|    |__|   \\___/ |_____||_____||___|___||__|__||____|\\___/  \\___|\n                                                                     \n")
	glg.Info("\"Σωτήρ\"")
	glg.Info("\"savior\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")

	baseConfig := configs.PtolemaiosConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	ptolemaiosConfig, ok := unparsedConfig.(*configs.PtolemaiosConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	srv := app.InitRoutes(*ptolemaiosConfig)

	glg.Infof("%s : %s", "running on port", port)
	err = http.ListenAndServe(port, srv)

	if err != nil {
		panic(err)
	}
}
