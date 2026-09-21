package listings

import (
	"database/sql"
	"log/slog"
	"net/http"
)

// Register wires up the listing feature (repository -> service -> handler)
// and attaches its routes to the given mux.
func Register(mux *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service, logger)

	mux.HandleFunc("GET /listings", handler.List)
	mux.HandleFunc("POST /listings", handler.Create)
	mux.HandleFunc("DELETE /listings/{id}", handler.Delete)
}
