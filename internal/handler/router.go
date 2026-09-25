package handler

import (
	"net/http"
)

func RegisterRoutes(linkH *LinkHandler, healthH *HealthHandler) http.Handler {
	mux := http.NewServeMux()

	// Links
	mux.HandleFunc("POST /api/v1/links", linkH.Create)
	mux.HandleFunc("GET /api/v1/links/{code}", linkH.GetByCode)
	mux.HandleFunc("DELETE /api/v1/links/{code}", linkH.Delete)

	// Health
	mux.HandleFunc("GET /health", healthH.Check)

	// Redirect
	mux.HandleFunc("GET /{code}", linkH.Redirect)

	// Middleware оборачивает mux
	var h http.Handler = mux
	// h = middleware.Logging(log)(h)
	// h = middleware.Recover(log)(h)

	return h
}
