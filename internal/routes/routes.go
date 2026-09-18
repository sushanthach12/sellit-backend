package routes

import (
	"database/sql"
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("GET /health", handlers.HealthCheck)

	mux.HandleFunc("GET /listings", handlers.List(db))
	mux.HandleFunc("DELETE /listings/{id}", handlers.DeleteListing(db))
}
