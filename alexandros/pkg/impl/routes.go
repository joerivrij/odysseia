package impl

import (
	"github.com/gorilla/mux"
	"github.com/lexiko/alexandros/pkg/config"
	"github.com/lexiko/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config config.AlexandrosConfig) *mux.Router {
	serveMux := mux.NewRouter()

	alexandrosHandler := AlexandrosHandler{Config: &config}

	serveMux.HandleFunc("/alexandros/v1/ping", middleware.Adapt(alexandrosHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/alexandros/v1/health", middleware.Adapt(alexandrosHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	return serveMux
}

