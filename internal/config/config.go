package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DSN         string
	Platform    string
	TokenSecret string
	PolkaKey    string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load .env file: %v", err)
	}
	return &Config{
		Port:        "8080",
		DSN:         os.Getenv("DB_URL"),
		Platform:    os.Getenv("PLATFORM"),
		TokenSecret: os.Getenv("TOKEN_SECRET"),
		PolkaKey:    os.Getenv("POLKA_KEY"),
	}, nil
}
