package handler

import (
	"log/slog"
	"net/http"
)

type HealthHandler struct {
	log *slog.Logger
}

func NewHealthHandler(log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		log: log,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}
