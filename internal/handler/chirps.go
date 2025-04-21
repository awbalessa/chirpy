package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errRes struct {
		Error string `json:"error"`
	}
	type cleanedResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := errRes{Error: "Something went wrong"}
		data, marshalErr := json.Marshal(errMsg)
		if marshalErr != nil {
			log.Printf("Error marshalling JSON: %v", err)
			return
		}
		w.WriteHeader(500)
		w.Write(data)
		return
	}

	if len(params.Body) > 140 {
		errMsg := errRes{Error: "Chirp is too long"}
		data, marshalErr := json.Marshal(errMsg)
		if marshalErr != nil {
			log.Printf("Error marshalling JSON: %v", err)
			return
		}
		w.WriteHeader(400)
		w.Write(data)
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
