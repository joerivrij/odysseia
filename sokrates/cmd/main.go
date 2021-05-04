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

	glg.Info("\n _______  _______  ___   _  ______    _______  _______  _______  _______ \n|       ||       ||   | | ||    _ |  |   _   ||       ||       ||       |\n|  _____||   _   ||   |_| ||   | ||  |  |_|  ||_     _||    ___||  _____|\n| |_____ |  | |  ||      _||   |_||_ |       |  |   |  |   |___ | |_____ \n|_____  ||  |_|  ||     |_ |    __  ||       |  |   |  |    ___||_____  |\n _____| ||       ||    _  ||   |  | ||   _   |  |   |  |   |___  _____| |\n|_______||_______||___| |_||___|  |_||__| |__|  |___|  |_______||_______|\n")
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
