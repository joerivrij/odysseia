package app

import (
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"os"
	"time"
)

func CreateHandler(config configs.PtolemaiosConfig) *PtolemaiosHandler {
	handler := PtolemaiosHandler{Config: &config, Duration: time.Second * 10}

	if config.RunOnce {
		go func() {
			jobExit := make(chan bool, 1)
			go handler.CheckForJobExit(jobExit)

			select {

			case <-jobExit:
				glg.Debug("exiting because of condition")
				os.Exit(0)
			}
		}()
	}

	return &handler
}
