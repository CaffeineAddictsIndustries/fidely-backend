package handler

import (
	"html/template"
	"net/http"
	"time"

	"fidely-backend/internal/config"
	"fidely-backend/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// WebHandler serves HTML pages for the admin interface.
type WebHandler struct {
	templates   *template.Template
	config      *config.Config
	authService *service.AdminAuthService
}

// NewWebHandler parses templates used by the admin web pages.
func NewWebHandler(cfg *config.Config, authService *service.AdminAuthService) (*WebHandler, error) {
	tmpl, err := template.ParseGlob("web/templates/*/*.html")
	if err != nil {
		return nil, err
	}

	return &WebHandler{templates: tmpl, config: cfg, authService: authService}, nil
}

// LoginPage renders the admin login page.
func (h *WebHandler) LoginPage(c echo.Context) error {
	return h.render(c, http.StatusOK, map[string]any{
		"Title": "Fidely Admin Login",
	})
}

// HandleLogin handles login POST requests.
func (h *WebHandler) HandleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	result, err := h.authService.Login(c.Request().Context(), username, password)
	if err != nil {
		status := http.StatusInternalServerError
		if isHTMX {
			status = http.StatusOK
		}

		return h.render(c, status, map[string]any{
			"Title":   "Fidely Admin Login",
			"Message": service.MessageLoginFailed,
			"Reason":  "unable to complete login right now",
			"Success": false,
		})
	}

	if !result.Success {
		status := http.StatusUnauthorized
		if isHTMX {
			status = http.StatusOK
		}

		return h.render(c, status, map[string]any{
			"Title":   "Fidely Admin Login",
			"Message": result.Message,
			"Reason":  service.ReasonInvalidCredentials,
			"Success": false,
		})
	}

	c.SetCookie(h.newSessionCookie(result.SessionToken, result.ExpiresAt))
	return h.render(c, http.StatusOK, map[string]any{
		"Title":   "Fidely Admin Login",
		"Message": result.Message,
		"Success": true,
	})
}

// HandleLogout revokes the active admin session.
func (h *WebHandler) HandleLogout(c echo.Context) error {
	sessionToken, ok := SessionTokenFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": service.MessageLogoutFailed,
			"reason":  service.ReasonInvalidSession,
		})
	}

	result, logoutErr := h.authService.Logout(c.Request().Context(), sessionToken)
	h.clearSessionCookie(c)
	if logoutErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": service.MessageLogoutFailed,
			"reason":  "unable to complete logout right now",
		})
	}

	status := http.StatusOK
	if !result.Success {
		status = http.StatusUnauthorized
	}

	return c.JSON(status, map[string]any{
		"success": result.Success,
		"message": result.Message,
		"reason":  result.Reason,
	})
}

// CurrentAdmin returns the authenticated admin bound to the current session cookie.
func (h *WebHandler) CurrentAdmin(c echo.Context) error {
	principal, ok := AuthenticatedPrincipalFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": service.MessageLoginFailed,
			"reason":  service.ReasonInvalidSession,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "authenticated",
		"admin": map[string]any{
			"id":         principal.AdminID,
			"admin_type": principal.AdminType,
			"username":   principal.Username,
			"role":       principal.Role,
			"store_id":   principal.StoreID,
		},
	})
}

// PlatformOnlyStatus verifies platform-admin authorization without introducing the later landing page flow.
func (h *WebHandler) PlatformOnlyStatus(c echo.Context) error {
	principal, ok := AuthenticatedPrincipalFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": service.MessageLoginFailed,
			"reason":  service.ReasonInvalidSession,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "platform access granted",
		"admin": map[string]any{
			"id":         principal.AdminID,
			"admin_type": principal.AdminType,
			"username":   principal.Username,
			"role":       principal.Role,
		},
	})
}

func (h *WebHandler) render(c echo.Context, code int, data map[string]any) error {
	if token, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string); ok && token != "" {
		data["CSRFToken"] = token
	}

	c.Response().WriteHeader(code)
	if c.Request().Header.Get("HX-Request") == "true" {
		return h.templates.ExecuteTemplate(c.Response().Writer, "content", data)
	}
	return h.templates.ExecuteTemplate(c.Response().Writer, "base", data)
}

func (h *WebHandler) newSessionCookie(token string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     h.config.AuthSessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.AuthCookieSecure,
		SameSite: toSameSiteMode(h.config.AuthCookieSameSite),
		Expires:  expiresAt,
	}
}

func (h *WebHandler) clearSessionCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     h.config.AuthSessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.AuthCookieSecure,
		SameSite: toSameSiteMode(h.config.AuthCookieSameSite),
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func toSameSiteMode(value string) http.SameSite {
	switch value {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
