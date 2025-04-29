package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/awbalessa/chirpy/internal/database"
	"github.com/awbalessa/chirpy/internal/middleware"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type APIConfig struct {
	Queries        *database.Queries
	FileServerHits atomic.Int32
	Platform       string
	TokenSecret    string
}

func (c *APIConfig) RespondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{Error: message}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return
	}
	w.Write(data)
}

func (cfg *APIConfig) WithMetrics(next http.Handler) http.Handler {
	return middleware.MetricsInc(&cfg.FileServerHits, next)
}
