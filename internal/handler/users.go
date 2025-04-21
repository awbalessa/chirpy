package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *APIConfig) handleUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		c.RespondWithError(w, 400, "Invalid request format")
		return
	}

	user, err := c.Queries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	res := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
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
