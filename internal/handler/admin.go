package handler

import (
	"fmt"
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

func (c *APIConfig) handleResetCounter(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	c.FileServerHits.Store(0)
}
