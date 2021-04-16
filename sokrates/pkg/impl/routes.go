package impl

import (
	"github.com/gorilla/mux"
	"sokrates/pkg/config"
	"sokrates/pkg/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config config.SokratesConfig) *mux.Router {
	serveMux := mux.NewRouter()

	sokratesHandler := SokratesHandler{Config: &config}

	serveMux.HandleFunc("/ping", middleware.Adapt(sokratesHandler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/api/v1/health", middleware.Adapt(sokratesHandler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/api/v1/chapters/{category}", middleware.Adapt(sokratesHandler.FindHighestChapter, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/api/v1/createQuestion", middleware.Adapt(sokratesHandler.CreateQuestion, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/api/v1/answer", middleware.Adapt(sokratesHandler.CheckAnswer, middleware.ValidateRestMethod("POST"), middleware.LogRequestDetails()))

	return serveMux
}