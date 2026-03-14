package config

import "testing"

func baseEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://fidely:fidely@localhost:5432/fidely?sslmode=disable")
	t.Setenv("SERVER_PORT", "8080")
}

func TestLoadRejectsInsecureSameSiteNone(t *testing.T) {
	baseEnv(t)
	t.Setenv("AUTH_SESSION_COOKIE_SAMESITE", "none")
	t.Setenv("AUTH_SESSION_COOKIE_SECURE", "false")
	t.Setenv("ENVIRONMENT", "development")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when sameSite is none and secure is false")
	}
}

func TestLoadRejectsProductionWithoutSecureCookie(t *testing.T) {
	baseEnv(t)
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("AUTH_SESSION_COOKIE_SECURE", "false")
	t.Setenv("AUTH_TOKEN_HASH_PEPPER", "strong-pepper")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when production secure cookie is false")
	}
}

func TestLoadRejectsProductionWithoutPepper(t *testing.T) {
	baseEnv(t)
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("AUTH_SESSION_COOKIE_SECURE", "true")
	t.Setenv("AUTH_TOKEN_HASH_PEPPER", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when production pepper is empty")
	}
}

func TestLoadAcceptsProductionSecureConfig(t *testing.T) {
	baseEnv(t)
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("AUTH_SESSION_COOKIE_SECURE", "true")
	t.Setenv("AUTH_SESSION_COOKIE_SAMESITE", "lax")
	t.Setenv("AUTH_TOKEN_HASH_PEPPER", "strong-pepper")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected valid production config, got error: %v", err)
	}
	if cfg.Environment != "production" {
		t.Fatalf("expected production environment, got %q", cfg.Environment)
	}
	if !cfg.AuthCookieSecure {
		t.Fatal("expected secure cookie true")
	}
	if cfg.AuthTokenHashPepper == "" {
		t.Fatal("expected non-empty pepper")
	}
}

func TestLoadNormalizesEnvironmentCase(t *testing.T) {
	baseEnv(t)
	t.Setenv("ENVIRONMENT", "Production")
	t.Setenv("AUTH_SESSION_COOKIE_SECURE", "true")
	t.Setenv("AUTH_SESSION_COOKIE_SAMESITE", "lax")
	t.Setenv("AUTH_TOKEN_HASH_PEPPER", "strong-pepper")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
	if cfg.Environment != "production" {
		t.Fatalf("expected normalized environment to be production, got %q", cfg.Environment)
	}
}

func TestLoadRejectsInvalidSameSiteValue(t *testing.T) {
	baseEnv(t)
	t.Setenv("AUTH_SESSION_COOKIE_SAMESITE", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when AUTH_SESSION_COOKIE_SAMESITE is invalid")
	}
}
