package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, cfg *APIConfig) {
	mux.HandleFunc("GET /api/healthz", cfg.handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleRequestCounter)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("POST /api/users", cfg.handleUsers)
	mux.HandleFunc("GET /api/chirps", cfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleGetChirpByID)
	mux.HandleFunc("POST /api/chirps", cfg.handlePostChirp)
}
