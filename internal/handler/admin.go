package handler

import (
	"fmt"
	"log"
	"net/http"
)

func (c *APIConfig) handleRequestCounter(w http.ResponseWriter, _ *http.Request) {
	html := `
	<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
	    <p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, html, c.FileServerHits.Load())
}

func (c *APIConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if c.Platform != "dev" {
		c.RespondWithError(w, 403, "Forbidden: reset endpoint is only available in development environments")
	}

	if err := c.Queries.ResetUsers(r.Context()); err != nil {
		log.Printf("Error resetting users: %v", err)
		c.RespondWithError(w, 500, "Internal server error")
	}

	w.WriteHeader(http.StatusOK)
}
