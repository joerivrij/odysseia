package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/periandros/app"
	"os"
	"strings"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PERIANDROS
	glg.Info("\n ____   ___  ____   ____   ____  ____   ___    ____   ___   _____\n|    \\ /  _]|    \\ |    | /    ||    \\ |   \\  |    \\ /   \\ / ___/\n|  o  )  [_ |  D  ) |  | |  o  ||  _  ||    \\ |  D  )     (   \\_ \n|   _/    _]|    /  |  | |     ||  |  ||  D  ||    /|  O  |\\__  |\n|  | |   [_ |    \\  |  | |  _  ||  |  ||     ||    \\|     |/  \\ |\n|  | |     ||  .  \\ |  | |  |  ||  |  ||     ||  .  \\     |\\    |\n|__| |_____||__|\\_||____||__|__||__|__||_____||__|\\_|\\___/  \\___|\n                                                                 \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"Περίανδρος δὲ ἦν Κυψέλου παῖς οὗτος ὁ τῷ Θρασυβούλῳ τὸ χρηστήριον μηνύσας· ἐτυράννευε δὲ ὁ Περίανδρος Κορίνθου\"")
	glg.Info("\"Periander, who disclosed the oracle's answer to Thrasybulus, was the son of Cypselus, and sovereign of Corinth\"")
	glg.Info(strings.Repeat("~", 37))

	glg.Debug("creating config")

	config := app.Get()

	healthy := config.CheckSolonHealth(120)
	if !healthy {
		glg.Fatal("death has found me")
	}

	created, err := config.CreateRole()
	if err != nil {
		glg.Fatal(err)
	}

	glg.Infof("user created was %t", created)
	os.Exit(0)
}
