package main

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/solon/app"
	"io/ioutil"
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

	env := os.Getenv("ENV")
	if env == "" {
		env = testingEnv
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=SOLON
	glg.Info("\n  _____  ___   _       ___   ____  \n / ___/ /   \\ | |     /   \\ |    \\ \n(   \\_ |     || |    |     ||  _  |\n \\__  ||  O  || |___ |  O  ||  |  |\n /  \\ ||     ||     ||     ||  |  |\n \\    ||     ||     ||     ||  |  |\n  \\___| \\___/ |_____| \\___/ |__|__|\n                                   \n")
	glg.Info("\"αὐτοὶ γὰρ οὐκ οἷοί τε ἦσαν αὐτὸ ποιῆσαι Ἀθηναῖοι: ὁρκίοισι γὰρ μεγάλοισι κατείχοντο δέκα ἔτεα χρήσεσθαι νόμοισι τοὺς ἄν σφι Σόλων θῆται.\"")
	glg.Info("\" since the Athenians themselves could not do that, for they were bound by solemn oaths to abide for ten years by whatever laws Solon should make.\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")


	var cert []byte
	var esClient *elasticsearch.Client
	if env != "TEST" {
		glg.Info("trying to read cert file from pod")
		cert, _ = ioutil.ReadFile("/app/config/certs/elastic-certificate.pem")
		es, err := elastic.CreateElasticClientFromEnvVariablesWithTLS(cert)
		if err != nil {
			glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
		}

		esClient = es
	} else {
		cert, _ = ioutil.ReadFile("/home/joerivrij/go/src/github.com/odysseia/solon/vault_config/elastic-certificate.pem")
		es, err := elastic.CreateElasticClientFromEnvVariablesWithTLS(cert)
		if err != nil {
			glg.Fatalf("Error creating ElasticClient shutting down: %s", err)
		}

		esClient = es
	}

	healthy, config := app.Get(200, esClient, cert, env)
	if !healthy {
		glg.Fatal("death has found me")
	}

	created := app.InitRoot(*config)
	glg.Info(created)
	srv := app.InitRoutes(*config)

	glg.Infof("%s : %s", "running on port", port)
	err := http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
