package app

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config configs.SokratesConfig) *mux.Router {
	serveMux := mux.NewRouter()

	sokratesHandler := SokratesHandler{Config: &config}

	serveMux.HandleFunc("/sokrates/v1/ping", middleware.Adapt(sokratesHandler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/health", middleware.Adapt(sokratesHandler.health, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/methods", middleware.Adapt(sokratesHandler.queryMethods, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/methods/{method}/categories", middleware.Adapt(sokratesHandler.queryCategories, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/methods/{method}/categories/{category}/chapters", middleware.Adapt(sokratesHandler.FindHighestChapter, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/createQuestion", middleware.Adapt(sokratesHandler.CreateQuestion, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/answer", middleware.Adapt(sokratesHandler.CheckAnswer, middleware.LogRequestDetails(), middleware.ValidateRestMethod("POST"), middleware.SetCorsHeaders()))

	return serveMux
}
