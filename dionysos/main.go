package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/dionysos/app"
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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=DIONYSOS
	glg.Info("\n ___    ____  ___   ____   __ __  _____  ___   _____\n|   \\  |    |/   \\ |    \\ |  |  |/ ___/ /   \\ / ___/\n|    \\  |  ||     ||  _  ||  |  (   \\_ |     (   \\_ \n|  D  | |  ||  O  ||  |  ||  ~  |\\__  ||  O  |\\__  |\n|     | |  ||     ||  |  ||___, |/  \\ ||     |/  \\ |\n|     | |  ||     ||  |  ||     |\\    ||     |\\    |\n|_____||____|\\___/ |__|__||____/  \\___| \\___/  \\___|\n                                                    \n")
	glg.Info("\"Γραμματική ἐστιν ἐμπειρία τῶν παρὰ ποιηταῖς τε καὶ συγγραφεῦσιν ὡς ἐπὶ τὸ πολὺ λεγομένων.’\"")
	glg.Info("\"Grammar is an experimental knowledge of the usages of language as generally current among poets and prose writers\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")

	baseConfig := configs.DionysosConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	dionysosConfig, ok := unparsedConfig.(*configs.DionysosConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	declensionConfig := app.QueryRuleSet(&dionysosConfig.ElasticClient, dionysosConfig.Index)
	dionysosConfig.DeclensionConfig = *declensionConfig

	srv := app.InitRoutes(*dionysosConfig)

	glg.Infof("%s : %s", "running on port", port)
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
