package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL         string
	Port                string
	JWTSecret           string
	GoogleClientID      string
	GoogleClientSecret  string
	FacebookClientID    string
	FacebookClientSecret string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		Port:                os.Getenv("PORT"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		FacebookClientID:    os.Getenv("FACEBOOK_CLIENT_ID"),
		FacebookClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
	}, nil
}