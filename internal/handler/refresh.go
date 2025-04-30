package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/awbalessa/chirpy/internal/auth"
)

func (c *APIConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	dbToken, err := c.Queries.GetRefreshToken(r.Context(), token)
	if err == sql.ErrNoRows {
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	} else if err != nil {
		log.Printf("Error getting refresh token: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	if !dbToken.ExpiresAt.Valid || time.Now().After(dbToken.ExpiresAt.Time) {
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	if dbToken.RevokedAt.Valid {
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	user, err := c.Queries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("Error getting user from refresh token: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	jwt, err := auth.MakeJWT(user.ID, c.TokenSecret, time.Hour)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	res := response{
		Token: jwt,
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling refresh token response: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
