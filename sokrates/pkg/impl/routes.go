package impl

import (
	"sokrates/pkg/config"
	"sokrates/pkg/middleware"
	"net/http"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(config config.SokratesConfig) *http.ServeMux {
	serveMux := http.NewServeMux()

	sokratesHandler := SokratesHandler{Config: &config}

	serveMux.HandleFunc("/ping", middleware.Adapt(sokratesHandler.PingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/api/v1/createDataSet", middleware.Adapt(sokratesHandler.CreateDocuments, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/api/v1/queryLogos", middleware.Adapt(sokratesHandler.QueryAllForIndex, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))

	return serveMux
}