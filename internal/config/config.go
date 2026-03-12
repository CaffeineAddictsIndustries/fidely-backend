package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application.
type Config struct {
	ServerPort          string
	DatabaseURL         string
	AuthSessionCookie   string
	AuthSessionTTL      time.Duration
	AuthCookieSecure    bool
	AuthCookieSameSite  string
	AuthTokenHashPepper string
}

// Load reads configuration from environment variables.
// Required env vars:
//   - DATABASE_URL: PostgreSQL connection string (e.g., postgres://user:pass@localhost:5432/fidely?sslmode=disable)
//   - SERVER_PORT: Port to run the HTTP server on (default: 8080)
//
// Optional auth env vars:
//   - AUTH_SESSION_COOKIE_NAME: Cookie name for admin sessions (default: fidely_admin_session)
//   - AUTH_SESSION_TTL: Session TTL using Go duration format (default: 12h)
//   - AUTH_SESSION_COOKIE_SECURE: Secure cookie flag (default: false)
//   - AUTH_SESSION_COOKIE_SAMESITE: lax|strict|none (default: lax)
//   - AUTH_TOKEN_HASH_PEPPER: Extra secret used in token hashing (default: empty)
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	cookieName := getenvOrDefault("AUTH_SESSION_COOKIE_NAME", "fidely_admin_session")

	ttlRaw := getenvOrDefault("AUTH_SESSION_TTL", "12h")
	ttl, err := time.ParseDuration(ttlRaw)
	if err != nil || ttl <= 0 {
		return nil, fmt.Errorf("AUTH_SESSION_TTL must be a positive Go duration, got %q", ttlRaw)
	}

	secureRaw := getenvOrDefault("AUTH_SESSION_COOKIE_SECURE", "false")
	secure, err := strconv.ParseBool(secureRaw)
	if err != nil {
		return nil, fmt.Errorf("AUTH_SESSION_COOKIE_SECURE must be true or false, got %q", secureRaw)
	}

	sameSite := strings.ToLower(getenvOrDefault("AUTH_SESSION_COOKIE_SAMESITE", "lax"))
	if sameSite != "lax" && sameSite != "strict" && sameSite != "none" {
		return nil, fmt.Errorf("AUTH_SESSION_COOKIE_SAMESITE must be one of lax|strict|none, got %q", sameSite)
	}

	return &Config{
		ServerPort:          port,
		DatabaseURL:         dbURL,
		AuthSessionCookie:   cookieName,
		AuthSessionTTL:      ttl,
		AuthCookieSecure:    secure,
		AuthCookieSameSite:  sameSite,
		AuthTokenHashPepper: os.Getenv("AUTH_TOKEN_HASH_PEPPER"),
	}, nil
}

func getenvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
