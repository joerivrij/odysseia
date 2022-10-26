package app

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia-greek/plato/aristoteles/configs"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config configs.DionysiosConfig) *mux.Router {
	serveMux := mux.NewRouter()

	dionysosHandler := DionysosHandler{Config: &config}

	serveMux.HandleFunc("/dionysios/v1/ping", middleware.Adapt(dionysosHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/dionysios/v1/health", middleware.Adapt(dionysosHandler.health, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/dionysios/v1/checkGrammar", middleware.Adapt(dionysosHandler.checkGrammar, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	return serveMux
}
