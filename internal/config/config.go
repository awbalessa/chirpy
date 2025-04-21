package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DSN  string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load .env file: %v", err)
	}
	return &Config{
		Port: "8080",
		DSN:  os.Getenv("postgres://azizalessa:@localhost:5432/chirpy"),
	}, nil
}
