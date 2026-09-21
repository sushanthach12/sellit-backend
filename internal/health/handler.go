package health

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Health check endpoint hit")

	w.Header().Set("Content-Type", "application/json") // Anything that is written after the WriteHeader call will be sent as the response body
	// so any headers modifications should be done before the WriteHeader call
	w.WriteHeader(http.StatusOK) // the order of the WriteHeader and Write methods is important; WriteHeader should be called before Write

	jsonResponse := `{"status": 200, "message": "Server Running Healthy..."}`

	w.Write([]byte(jsonResponse))
}
