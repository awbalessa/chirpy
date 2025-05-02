package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/awbalessa/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *APIConfig) handleUpgradeToRed(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	if key != c.PolkaKey {
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding request body: %v", err)
		c.RespondWithError(w, 400, "Invalid request format")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = c.Queries.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err == sql.ErrNoRows {
		c.RespondWithError(w, 404, "User not found")
		return
	} else if err != nil {
		log.Printf("Error upgrading to chirpy red: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.WriteHeader(204)
}
