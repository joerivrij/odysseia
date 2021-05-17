package main

import (
	"github.com/ianschenck/envflag"
	"github.com/kpango/glg"
	"github.com/lexiko/alexandros/pkg/config"
	"github.com/lexiko/alexandros/pkg/impl"
	"net/http"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	port := envflag.String("PORT", ":5000", "port")

	envflag.Parse()

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HERODOTOS
	glg.Info("\n  ____  _        ___  __ __   ____  ____   ___    ____   ___   _____\n /    || |      /  _]|  |  | /    ||    \\ |   \\  |    \\ /   \\ / ___/\n|  o  || |     /  [_ |  |  ||  o  ||  _  ||    \\ |  D  )     (   \\_ \n|     || |___ |    _]|_   _||     ||  |  ||  D  ||    /|  O  |\\__  |\n|  _  ||     ||   [_ |     ||  _  ||  |  ||     ||    \\|     |/  \\ |\n|  |  ||     ||     ||  |  ||  |  ||  |  ||     ||  .  \\     |\\    |\n|__|__||_____||_____||__|__||__|__||__|__||_____||__|\\_|\\___/  \\___|\n                                                                    \n")
	glg.Info("\"ὅτι τοῦ κρατεῖν πέρας ἡμῖν ἐστι τὸ μὴ ταὐτὰ ποιεῖν τοῖς κεκρατημένοις;’\"")
	glg.Info("\"Know ye not,’ said he, ‘that the end and object of conquest is to avoid doing the same thing as the conquered?\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")

	config := config.Get()

	srv := impl.InitRoutes(*config)

	glg.Infof("%s : %s", "running on port", *port)
	err := http.ListenAndServe(*port, srv)
	if err != nil {
		panic(err)
	}
}
