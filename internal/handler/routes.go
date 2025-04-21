package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, cfg *APIConfig) {
	mux.HandleFunc("GET /api/healthz", cfg.handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleRequestCounter)
	mux.HandleFunc("POST /admin/reset", cfg.handleResetCounter)
	mux.HandleFunc("POST /api/validate_chirp", cfg.handleValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.handleUsers)
}
