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

	type errMsg struct {
		Error string `json:"error"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		return
	}

	user, err := c.Queries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
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
		log.Printf("Error marshalling JSON: %v", err)
	}
	w.WriteHeader(201)
	w.Write(data)
	return
}
