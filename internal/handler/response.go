package handler

import (
	"encoding/json"
	"net/http"
)

func (h *LinkHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("encode response failed", "error", err)
	}
}

func (h *LinkHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, ErrorResponse{Error: message})
}
