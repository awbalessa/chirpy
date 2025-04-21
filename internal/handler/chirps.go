package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (c *APIConfig) handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type cleanedResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		c.RespondWithError(w, 400, "Invalid request format")
		return
	}

	if len(params.Body) > 140 {
		c.RespondWithError(w, 400, "Chirp is too long")
		return
	} else {
		resBodyArr := strings.Fields(params.Body)
		for i := range resBodyArr {
			if strings.ToLower(resBodyArr[i]) == "kerfuffle" || strings.ToLower(resBodyArr[i]) == "sharbert" || strings.ToLower(resBodyArr[i]) == "fornax" {
				resBodyArr[i] = "****"
			}
		}

		cleanedRes := cleanedResponse{CleanedBody: strings.Join(resBodyArr, " ")}
		data, err := json.Marshal(&cleanedRes)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
			return
		}
		w.WriteHeader(200)
		w.Write(data)
		return
	}
}
