package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/thrasyboulos/app"
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
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=THRASYBOULOS
	glg.Info("\n ______  __ __  ____    ____  _____ __ __  ____    ___   __ __  _       ___   _____\n|      ||  |  ||    \\  /    |/ ___/|  |  ||    \\  /   \\ |  |  || |     /   \\ / ___/\n|      ||  |  ||  D  )|  o  (   \\_ |  |  ||  o  )|     ||  |  || |    |     (   \\_ \n|_|  |_||  _  ||    / |     |\\__  ||  ~  ||     ||  O  ||  |  || |___ |  O  |\\__  |\n  |  |  |  |  ||    \\ |  _  |/  \\ ||___, ||  O  ||     ||  :  ||     ||     |/  \\ |\n  |  |  |  |  ||  .  \\|  |  |\\    ||     ||     ||     ||     ||     ||     |\\    |\n  |__|  |__|__||__|\\_||__|__| \\___||____/ |_____| \\___/  \\__,_||_____| \\___/  \\___|\n                                                                                   \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"πέμψας γὰρ παρὰ Θρασύβουλον κήρυκα ἐπυνθάνετο ὅντινα ἂν τρόπον ἀσφαλέστατον καταστησάμενος τῶν πρηγμάτων κάλλιστα τὴν πόλιν ἐπιτροπεύοι.\"")
	glg.Info("\"He had sent a herald to Thrasybulus and inquired in what way he would best and most safely govern his city. \"")
	glg.Info(strings.Repeat("~", 37))

	glg.Debug("creating config")

	baseConfig := configs.ThrasyboulosConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	config, ok := unparsedConfig.(*configs.ThrasyboulosConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	done := make(chan bool)
	handler := app.ThrasyboulosHandler{Config: config}

	go func() {
		handler.WaitForJobsToFinish(done)
	}()

	select {

	case <-done:
		glg.Infof("%s job finished", config.Job)
		os.Exit(0)
	}

}
