package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/awbalessa/chirpy/internal/config"
	"github.com/awbalessa/chirpy/internal/database"
	"github.com/awbalessa/chirpy/internal/handler"
	_ "github.com/lib/pq"
)

func StaticHandler() http.Handler {
	dir := filepath.Join("web", "app")
	fs := http.FileServer(http.Dir(dir))
	return http.StripPrefix("/app/", fs)
}

func Run(cfg config.Config) error {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return fmt.Errorf("Error opening DB: %v", err)
	}
	defer db.Close()

	apiCfg := &handler.APIConfig{
		Queries:     database.New(db),
		Platform:    cfg.Platform,
		TokenSecret: cfg.TokenSecret,
		PolkaKey:    cfg.PolkaKey,
	}

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, apiCfg)
	mux.Handle("/app/", apiCfg.WithMetrics(StaticHandler()))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	fmt.Printf("Listening on http://localhost:%s\n", cfg.Port)
	return server.ListenAndServe()
}
