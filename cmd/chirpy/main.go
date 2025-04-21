package main

import (
	"log"

	"github.com/awbalessa/chirpy/internal/config"
	"github.com/awbalessa/chirpy/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	err = server.Run(*cfg)
	if err != nil {
		log.Fatal(err)
	}
}
