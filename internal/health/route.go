package health

import (
	"log/slog"
	"net/http"
)

func Register(mux *http.ServeMux, logger *slog.Logger) {
	handler := NewHandler(logger)
	mux.HandleFunc("GET /health", handler.HealthCheck)
}
