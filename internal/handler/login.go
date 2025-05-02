package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/awbalessa/chirpy/internal/auth"
	"github.com/awbalessa/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *APIConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		c.RespondWithError(w, 400, "Invalid request format")
		return
	}

	user, err := c.Queries.GetUserByEmail(r.Context(), params.Email)
	if err == sql.ErrNoRows {
		c.RespondWithError(w, 401, "Incorrect email or password")
		return
	} else if err != nil {
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	if err = auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		log.Print(err)
		c.RespondWithError(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, c.TokenSecret, time.Hour)
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Print(err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	refreshParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}
	_, err = c.Queries.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	res := response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed.Bool,
		Token:        token,
		RefreshToken: refreshToken,
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
	return
}
