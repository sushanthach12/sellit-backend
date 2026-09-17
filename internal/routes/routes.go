package routes

import (
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", handlers.HealthCheck)
}
