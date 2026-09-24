package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	ServerPort    string
	Environment   string
	AllowedOrigin string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		ServerPort:    os.Getenv("SERVER_PORT"),
		Environment:   os.Getenv("ENVIRONMENT"),
		AllowedOrigin: os.Getenv("ALLOWED_ORIGIN"),
	}
	if cfg.AllowedOrigin == "" {
		cfg.AllowedOrigin = "http://localhost:3000"
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABSE_URL is required but was not set")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required but was not set")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters long, got %d", len(cfg.JWTSecret))
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}
	return cfg, nil
}
