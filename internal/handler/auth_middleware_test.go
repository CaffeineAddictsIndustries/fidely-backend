package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fidely-backend/internal/auth"
	"fidely-backend/internal/config"
	"fidely-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type authMiddlewareRepo struct {
	byUsername map[string]*service.AdminPrincipal
	byKey      map[string]*service.AdminPrincipal
	sessions   map[string]auth.AdminSession
	nextID     int
}

func newAuthMiddlewareRepo() *authMiddlewareRepo {
	return &authMiddlewareRepo{
		byUsername: make(map[string]*service.AdminPrincipal),
		byKey:      make(map[string]*service.AdminPrincipal),
		sessions:   make(map[string]auth.AdminSession),
		nextID:     1,
	}
}

func middlewarePrincipalKey(adminType auth.AdminType, adminID int) string {
	return fmt.Sprintf("%s:%d", adminType, adminID)
}

func (repo *authMiddlewareRepo) FindByUsername(_ context.Context, username string) (*service.AdminPrincipal, error) {
	return repo.byUsername[username], nil
}

func (repo *authMiddlewareRepo) FindByTypeAndID(_ context.Context, adminType auth.AdminType, adminID int) (*service.AdminPrincipal, error) {
	return repo.byKey[middlewarePrincipalKey(adminType, adminID)], nil
}

func (repo *authMiddlewareRepo) CreateSession(_ context.Context, session auth.AdminSession) (auth.AdminSession, error) {
	session.ID = repo.nextID
	repo.nextID++
	repo.sessions[session.SessionTokenHash] = session
	return session, nil
}

func (repo *authMiddlewareRepo) GetSessionByTokenHash(_ context.Context, tokenHash string) (*auth.AdminSession, error) {
	session, ok := repo.sessions[tokenHash]
	if !ok {
		return nil, nil
	}
	copySession := session
	return &copySession, nil
}

func (repo *authMiddlewareRepo) UpdateSession(_ context.Context, session auth.AdminSession) error {
	repo.sessions[session.SessionTokenHash] = session
	return nil
}

func buildAuthenticatedService(t *testing.T, principal *service.AdminPrincipal, rawPassword string) (*service.AdminAuthService, *config.Config, string) {
	t.Helper()

	repo := newAuthMiddlewareRepo()
	passwords := auth.NewDefaultPasswordManager()
	sessions := auth.NewSessionManager("pepper")

	hash, err := passwords.Hash(rawPassword)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	principal.PasswordHash = hash
	repo.byUsername[principal.Username] = principal
	repo.byKey[middlewarePrincipalKey(principal.AdminType, principal.AdminID)] = principal

	authService, err := service.NewAdminAuthService(repo, passwords, sessions, time.Hour)
	if err != nil {
		t.Fatalf("failed to create auth service: %v", err)
	}

	loginResult, err := authService.Login(context.Background(), principal.Username, rawPassword)
	if err != nil {
		t.Fatalf("failed to login test principal: %v", err)
	}
	if !loginResult.Success {
		t.Fatalf("expected login success, got: %+v", loginResult)
	}

	cfg := &config.Config{AuthSessionCookie: "fidely_admin_session"}
	return authService, cfg, loginResult.SessionToken
}

func TestRequireAuthenticatedAdminRejectsMissingCookie(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &service.AdminPrincipal{AdminType: auth.AdminTypeStoreAdmin, AdminID: 1, Username: "alice", Role: "manager"}
	authService, cfg, _ := buildAuthenticatedService(t, principal, "Password123!")
	middleware := NewAuthMiddleware(cfg, authService)

	handler := middleware.RequireAuthenticatedAdmin()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if err := handler(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAuthenticatedAdminInjectsPrincipal(t *testing.T) {
	e := echo.New()
	principal := &service.AdminPrincipal{AdminType: auth.AdminTypeStoreAdmin, AdminID: 2, Username: "alice", Role: "manager"}
	authService, cfg, sessionToken := buildAuthenticatedService(t, principal, "Password123!")
	middleware := NewAuthMiddleware(cfg, authService)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: cfg.AuthSessionCookie, Value: sessionToken})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware.RequireAuthenticatedAdmin()(func(c echo.Context) error {
		resolvedPrincipal, ok := AuthenticatedPrincipalFromContext(c)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.String(http.StatusOK, resolvedPrincipal.Username)
	})

	if err := handler(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "alice" {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestRequireAdminTypesRejectsInsufficientPrivileges(t *testing.T) {
	e := echo.New()
	principal := &service.AdminPrincipal{AdminType: auth.AdminTypeStoreAdmin, AdminID: 3, Username: "bob", Role: "manager"}
	authService, cfg, sessionToken := buildAuthenticatedService(t, principal, "Password123!")
	middleware := NewAuthMiddleware(cfg, authService)

	req := httptest.NewRequest(http.MethodGet, "/admin/platform/status", nil)
	req.AddCookie(&http.Cookie{Name: cfg.AuthSessionCookie, Value: sessionToken})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware.RequireAuthenticatedAdmin()(middleware.RequireAdminTypes(auth.AdminTypeFidelyAdmin)(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))

	if err := handler(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestRequireAdminTypesAllowsFidelyAdmin(t *testing.T) {
	e := echo.New()
	principal := &service.AdminPrincipal{AdminType: auth.AdminTypeFidelyAdmin, AdminID: 4, Username: "platform", Role: "1"}
	authService, cfg, sessionToken := buildAuthenticatedService(t, principal, "Password123!")
	middleware := NewAuthMiddleware(cfg, authService)

	req := httptest.NewRequest(http.MethodGet, "/admin/platform/status", nil)
	req.AddCookie(&http.Cookie{Name: cfg.AuthSessionCookie, Value: sessionToken})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware.RequireAuthenticatedAdmin()(middleware.RequireAdminTypes(auth.AdminTypeFidelyAdmin)(func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	}))

	if err := handler(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
