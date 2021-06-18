package impl

import (
	"github.com/gorilla/mux"
	"github.com/odysseia/sokrates/pkg/config"
	"github.com/odysseia/sokrates/pkg/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config config.SokratesConfig) *mux.Router {
	serveMux := mux.NewRouter()

	sokratesHandler := SokratesHandler{Config: &config}

	serveMux.HandleFunc("/sokrates/v1/ping", middleware.Adapt(sokratesHandler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/health", middleware.Adapt(sokratesHandler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/chapters/{category}", middleware.Adapt(sokratesHandler.FindHighestChapter, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/createQuestion", middleware.Adapt(sokratesHandler.CreateQuestion, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/answer", middleware.Adapt(sokratesHandler.CheckAnswer, middleware.LogRequestDetails(), middleware.ValidateRestMethod("POST"), middleware.SetCorsHeaders()))

	return serveMux
}
