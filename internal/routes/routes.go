package routes

import (
	"database/sql"
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("GET /health", handlers.HealthCheck)

	listingHandler := handlers.NewListingHandler(db)
	mux.HandleFunc("GET /listings", listingHandler.List)
	mux.HandleFunc("DELETE /listings/{id}", listingHandler.Delete)
}
