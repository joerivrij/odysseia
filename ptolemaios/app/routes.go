package app

import (
	"github.com/gorilla/mux"
	"github.com/odysseia/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config PtolemaiosConfig) *mux.Router {
	serveMux := mux.NewRouter()

	handler := PtolemaiosHandler{Config: &config}

	serveMux.HandleFunc("/ptolemaios/v1/ping", middleware.Adapt(handler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/ptolemaios/v1/secret", middleware.Adapt(handler.GetSecretFromVault, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	if config.IsPartOfJob {
		go handler.CheckForJobExit()
	}

	return serveMux
}
