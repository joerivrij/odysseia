package impl

import (
	"github.com/gorilla/mux"
	"github.com/lexiko/herodotos/pkg/config"
	"github.com/lexiko/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config config.HerodotosConfig) *mux.Router {
	serveMux := mux.NewRouter()

	herodotosHandler := HerodotosHandler{Config: &config}

	serveMux.HandleFunc("/herodotos/v1/ping", middleware.Adapt(herodotosHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/herodotos/v1/createQuestion", middleware.Adapt(herodotosHandler.createQuestion, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/herodotos/v1/checkSentence", middleware.Adapt(herodotosHandler.checkSentence, middleware.ValidateRestMethod("POST"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))


	return serveMux
}