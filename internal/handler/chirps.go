package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/awbalessa/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *APIConfig) handlePostChirp(w http.ResponseWriter, r *http.Request) {
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

func (c *APIConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps, err := c.Queries.GetChirpsOldestFirst(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	responseChirps := make([]chirp, len(chirps))

	for i, dbChirp := range chirps {
		responseChirps[i] = chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	data, err := json.Marshal(responseChirps)
	if err != nil {
		log.Printf("Error marshalling response chirps: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
	return
}

func (c *APIConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	inputID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		c.RespondWithError(w, 400, "Invalid chirp ID format")
		return
	}

	dbChirp, err := c.Queries.GetChirpByID(r.Context(), inputID)
	if err == sql.ErrNoRows {
		log.Printf("Chirp not found: %v", err)
		c.RespondWithError(w, 404, "Chirp not found")
		return
	} else if err != nil {
		log.Printf("Error getting chirp by ID: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	resChirp := chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	data, err := json.Marshal(resChirp)
	if err != nil {
		log.Printf("Error marshalling chirp: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
	return
}
