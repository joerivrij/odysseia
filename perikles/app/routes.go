package app

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia-greek/plato/aristoteles/configs"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config configs.PeriklesConfig) *mux.Router {
	serveMux := mux.NewRouter()

	periklesHandler := PeriklesHandler{Config: &config}

	serveMux.HandleFunc("/perikles/v1/ping", middleware.Adapt(periklesHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/perikles/v1/validate", middleware.Adapt(periklesHandler.validate, middleware.ValidateRestMethod("POST"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	go periklesHandler.loopForMappingUpdates()

	return serveMux
}
