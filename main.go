package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := &apiConfig{}
	mux := http.NewServeMux()
	handleRoot := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfg.middlewareMetricsInc(handleRoot))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleRequestCounter)
	mux.HandleFunc("POST /admin/reset", cfg.handleResetCounter)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server shut down gracefully")
		} else {
			log.Fatalf("Unexpected server error: %v", err)
		}
	}
}

func handleReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (c *apiConfig) handleRequestCounter(w http.ResponseWriter, _ *http.Request) {
	html := `
	<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
	    <p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, html, c.fileserverHits.Load())
}

func (c *apiConfig) handleResetCounter(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	c.fileserverHits.Store(0)
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errRes struct {
		Error string `json:"error"`
	}
	type valid struct {
		Valid bool `json:"valid"`
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
