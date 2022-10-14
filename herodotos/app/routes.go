package app

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia/aristoteles/configs"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config configs.HerodotosConfig) *mux.Router {
	serveMux := mux.NewRouter()

	herodotosHandler := HerodotosHandler{Config: &config}

	serveMux.HandleFunc("/herodotos/v1/ping", middleware.Adapt(herodotosHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/herodotos/v1/health", middleware.Adapt(herodotosHandler.health, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/herodotos/v1/createQuestion", middleware.Adapt(herodotosHandler.createQuestion, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/herodotos/v1/authors", middleware.Adapt(herodotosHandler.queryAuthors, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/herodotos/v1/authors/{author}/books", middleware.Adapt(herodotosHandler.queryBooks, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/herodotos/v1/checkSentence", middleware.Adapt(herodotosHandler.checkSentence, middleware.ValidateRestMethod("POST"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	return serveMux
}
