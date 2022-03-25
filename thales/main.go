package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/queue"
	"github.com/odysseia/thales/app"
	"os"
	"strings"
	"time"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=THALES
	glg.Info("\n ______  __ __   ____  _        ___  _____\n|      ||  |  | /    || |      /  _]/ ___/\n|      ||  |  ||  o  || |     /  [_(   \\_ \n|_|  |_||  _  ||     || |___ |    _]\\__  |\n  |  |  |  |  ||  _  ||     ||   [_ /  \\ |\n  |  |  |  |  ||  |  ||     ||     |\\    |\n  |__|  |__|__||__|__||_____||_____| \\___|\n                                          \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"Μέγιστον τόπος· ἄπαντα γὰρ χωρεῖ.\"")
	glg.Info("\"The greatest is space, for it holds all things\"")
	glg.Info(strings.Repeat("~", 37))

	baseConfig := configs.ThalesConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	thalesConfig, ok := unparsedConfig.(*configs.ThalesConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	defer func(Queue queue.Client) {
		err := Queue.Close()
		if err != nil {
			glg.Error(err)
		}
	}(thalesConfig.Queue)

	duration := 10000 * time.Millisecond
	handler := app.ThalesHandler{Config: thalesConfig, QueueEmptyDuration: duration}

	handler.HandleQueue()

	glg.Info("queue empty - all messages handled")

	os.Exit(0)
}
