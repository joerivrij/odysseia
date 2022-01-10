package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/solon/app"
	"net/http"
	"os"
)

const standardPort = ":5000"
const testingEnv = "TEST"

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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=SOLON
	glg.Info("\n  _____  ___   _       ___   ____  \n / ___/ /   \\ | |     /   \\ |    \\ \n(   \\_ |     || |    |     ||  _  |\n \\__  ||  O  || |___ |  O  ||  |  |\n /  \\ ||     ||     ||     ||  |  |\n \\    ||     ||     ||     ||  |  |\n  \\___| \\___/ |_____| \\___/ |__|__|\n                                   \n")
	glg.Info("\"αὐτοὶ γὰρ οὐκ οἷοί τε ἦσαν αὐτὸ ποιῆσαι Ἀθηναῖοι: ὁρκίοισι γὰρ μεγάλοισι κατείχοντο δέκα ἔτεα χρήσεσθαι νόμοισι τοὺς ἄν σφι Σόλων θῆται.\"")
	glg.Info("\"since the Athenians themselves could not do that, for they were bound by solemn oaths to abide for ten years by whatever laws Solon should make.\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")

	baseConfig := configs.SolonConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	solonConfig, ok := unparsedConfig.(*configs.SolonConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	srv := app.InitRoutes(*solonConfig)

	glg.Infof("%s : %s", "running on port", port)
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
