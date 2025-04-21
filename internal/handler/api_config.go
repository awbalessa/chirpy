package handler

import (
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
}

func (cfg *APIConfig) WithMetrics(next http.Handler) http.Handler {
	return middleware.MetricsInc(&cfg.FileServerHits, next)
}
