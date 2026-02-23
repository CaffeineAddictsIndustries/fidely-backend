package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application.
type Config struct {
	ServerPort  string
	DatabaseURL string
}

// Load reads configuration from environment variables.
// Required env vars:
//   - DATABASE_URL: PostgreSQL connection string (e.g., postgres://user:pass@localhost:5432/fidely?sslmode=disable)
//   - SERVER_PORT: Port to run the HTTP server on (default: 8080)
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		ServerPort:  port,
		DatabaseURL: dbURL,
	}, nil
}
