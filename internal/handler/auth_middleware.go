package handler

import (
	"net/http"

	"fidely-backend/internal/auth"
	"fidely-backend/internal/config"
	"fidely-backend/internal/service"

	"github.com/labstack/echo/v4"
)

const (
	adminPrincipalContextKey = "auth.adminPrincipal"
	sessionTokenContextKey   = "auth.sessionToken"
)

// AuthMiddleware resolves and authorizes authenticated admin principals.
type AuthMiddleware struct {
	config      *config.Config
	authService *service.AdminAuthService
}

func NewAuthMiddleware(cfg *config.Config, authService *service.AdminAuthService) *AuthMiddleware {
	return &AuthMiddleware{config: cfg, authService: authService}
}

func (middleware *AuthMiddleware) RequireAuthenticatedAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(middleware.config.AuthSessionCookie)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": service.MessageLoginFailed,
					"reason":  service.ReasonInvalidSession,
				})
			}

			principal, err := middleware.authService.Authenticate(c.Request().Context(), cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": service.MessageLoginFailed,
					"reason":  service.ReasonInvalidSession,
				})
			}

			c.Set(adminPrincipalContextKey, principal)
			c.Set(sessionTokenContextKey, cookie.Value)
			return next(c)
		}
	}
}

func (middleware *AuthMiddleware) RequireAdminTypes(allowedTypes ...auth.AdminType) echo.MiddlewareFunc {
	allowed := make(map[auth.AdminType]struct{}, len(allowedTypes))
	for _, adminType := range allowedTypes {
		allowed[adminType] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			principal, ok := AuthenticatedPrincipalFromContext(c)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": service.MessageLoginFailed,
					"reason":  service.ReasonInvalidSession,
				})
			}

			if _, exists := allowed[principal.AdminType]; !exists {
				return c.JSON(http.StatusForbidden, map[string]any{
					"success": false,
					"message": service.MessageAccessDenied,
					"reason":  service.ReasonInsufficientPrivileges,
				})
			}

			return next(c)
		}
	}
}

func AuthenticatedPrincipalFromContext(c echo.Context) (*service.AdminPrincipal, bool) {
	principal, ok := c.Get(adminPrincipalContextKey).(*service.AdminPrincipal)
	return principal, ok && principal != nil
}

func SessionTokenFromContext(c echo.Context) (string, bool) {
	token, ok := c.Get(sessionTokenContextKey).(string)
	return token, ok && token != ""
}
