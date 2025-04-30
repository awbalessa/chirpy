package handler

import (
	"log"
	"net/http"

	"github.com/awbalessa/chirpy/internal/auth"
)

func (c *APIConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	if err = c.Queries.RevokeToken(r.Context(), token); err != nil {
		log.Printf("Error revoking token: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.WriteHeader(204)
}
