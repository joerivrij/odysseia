package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/perikles/app"
	"strings"
)

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PERIKLES
	glg.Info("\n ____   ___  ____   ____  __  _  _        ___  _____\n|    \\ /  _]|    \\ |    ||  |/ ]| |      /  _]/ ___/\n|  o  )  [_ |  D  ) |  | |  ' / | |     /  [_(   \\_ \n|   _/    _]|    /  |  | |    \\ | |___ |    _]\\__  |\n|  | |   [_ |    \\  |  | |     ||     ||   [_ /  \\ |\n|  | |     ||  .  \\ |  | |  .  ||     ||     |\\    |\n|__| |_____||__|\\_||____||__|\\_||_____||_____| \\___|\n                                                    \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"τόν γε σοφώτατον οὐχ ἁμαρτήσεται σύμβουλον ἀναμείνας χρόνον.\"")
	glg.Info("\"he would yet do full well to wait for that wisest of all counsellors, Time.\"")
	glg.Info(strings.Repeat("~", 37))

	glg.Debug("creating config")

	baseConfig := configs.PeriklesConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	periklesConfig, ok := unparsedConfig.(*configs.PeriklesConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	handler := app.PeriklesHandler{Config: periklesConfig}

	glg.Info("init for CA started...")
	err = handler.Config.Cert.InitCa()
	if err != nil {
		glg.Fatal(err)
	}

	glg.Info("CA created starting app")

	handler.Flow()

}
