package routes

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/health"
	"github.com/sushanthach12/sellit-backend/internal/listings"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	health.Register(mux, logger)

	listings.Register(mux, db, logger)
}
