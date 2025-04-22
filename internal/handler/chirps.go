package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/awbalessa/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *APIConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	var params reqParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		c.RespondWithError(w, 400, "Invalid request format")
		return
	}

	if len(params.Body) > 140 {
		c.RespondWithError(w, 400, "Chirp is too long")
		return
	}
	bodyArr := strings.Fields(params.Body)
	for i := range bodyArr {
		switch strings.ToLower(bodyArr[i]) {
		case "kerfuffle", "sharbert", "fornax":
			bodyArr[i] = "****"
		}
	}

	chirpParams := database.CreateChirpParams{
		Body:   strings.Join(bodyArr, " "),
		UserID: params.UserID,
	}
	chirp, err := c.Queries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	res := response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	data, err := json.Marshal(res)
	if err != nil {
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
	return
}
