package main

import (
	"github.com/ianschenck/envflag"
	"github.com/kpango/glg"
	"github.com/lexiko/herodotos/pkg/config"
	"github.com/lexiko/herodotos/pkg/impl"
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
	glg.Info("\n __ __    ___  ____   ___   ___     ___   ______   ___   _____\n|  |  |  /  _]|    \\ /   \\ |   \\   /   \\ |      | /   \\ / ___/\n|  |  | /  [_ |  D  )     ||    \\ |     ||      ||     (   \\_ \n|  _  ||    _]|    /|  O  ||  D  ||  O  ||_|  |_||  O  |\\__  |\n|  |  ||   [_ |    \\|     ||     ||     |  |  |  |     |/  \\ |\n|  |  ||     ||  .  \\     ||     ||     |  |  |  |     |\\    |\n|__|__||_____||__|\\_|\\___/ |_____| \\___/   |__|   \\___/  \\___|\n                                                              \n")
	glg.Info("\"Ἡροδότου Ἁλικαρνησσέος ἱστορίης ἀπόδεξις ἥδε\"")
	glg.Info("\"This is the display of the inquiry of Herodotos of Halikarnassos\"")
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
