package routes

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	mux.HandleFunc("GET /health", handlers.HealthCheck)

	listingHandler := handlers.NewListingHandler(db, logger)
	mux.HandleFunc("GET /listings", listingHandler.List)
	mux.HandleFunc("POST /listings", listingHandler.Create)
	mux.HandleFunc("DELETE /listings/{id}", listingHandler.Delete)
}
