package main

import (
	"github.com/ianschenck/envflag"
	"github.com/kpango/glg"
	"net/http"
	"sokrates/pkg/config"
	"sokrates/pkg/impl"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	port := envflag.String("PORT",":5000", "port")

	envflag.Parse()

	glg.Info("welcome to harpokratos")
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
