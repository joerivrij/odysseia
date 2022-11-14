package app

import (
	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/middleware"
	"os"
	"time"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config configs.PtolemaiosConfig) *mux.Router {
	serveMux := mux.NewRouter()

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

	serveMux.HandleFunc("/ptolemaios/v1/ping", middleware.Adapt(handler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/ptolemaios/v1/health", middleware.Adapt(handler.Health, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/ptolemaios/v1/secret", middleware.Adapt(handler.GetSecretFromVault, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	return serveMux
}

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
