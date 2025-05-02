package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/awbalessa/chirpy/internal/auth"
	"github.com/awbalessa/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *APIConfig) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		c.RespondWithError(w, 400, "Invalid request format")
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	dbParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPwd,
	}

	user, err := c.Queries.CreateUser(r.Context(), dbParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	res := response{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
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

func (c *APIConfig) handlePutUsers(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	userID, err := auth.ValidateJWT(token, c.TokenSecret)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 401, "Unauthorized to access resource")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&params); err != nil {
		log.Printf("Error decoding: %v", err)
		c.RespondWithError(w, 400, "Invalid request format")
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	dbParams := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashed,
	}

	newUser, err := c.Queries.UpdateUser(r.Context(), dbParams)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	res := response{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed.Bool,
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling new user: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
