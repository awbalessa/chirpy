package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, cfg *APIConfig) {
	mux.HandleFunc("GET /api/healthz", cfg.handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleRequestCounter)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("POST /api/users", cfg.handlePostUsers)
	mux.HandleFunc("PUT /api/users", cfg.handlePutUsers)
	mux.HandleFunc("GET /api/chirps", cfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handleDeleteChirp)
	mux.HandleFunc("POST /api/chirps", cfg.handlePostChirp)
	mux.HandleFunc("POST /api/login", cfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handleRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.handleUpgradeToRed)
}
