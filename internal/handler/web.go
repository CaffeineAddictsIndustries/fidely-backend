package handler

import (
    "html/template"
    "net/http"

    "github.com/labstack/echo/v4"
)

// WebHandler serves HTML pages for the admin interface.
type WebHandler struct {
    templates *template.Template
}

// NewWebHandler parses templates used by the admin web pages.
func NewWebHandler() (*WebHandler, error) {
    tmpl, err := template.ParseGlob("web/templates/*/*.html")
    if err != nil {
        return nil, err
    }

    return &WebHandler{templates: tmpl}, nil
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

    // Placeholder validation until auth service is implemented.
    if username == "" || password == "" {
        return h.render(c, http.StatusBadRequest, map[string]any{
            "Title": "Fidely Admin Login",
            "Error": "Please enter both username and password.",
        })
    }

    c.Response().Header().Set("HX-Redirect", "/dashboard")
    return c.NoContent(http.StatusOK)
}

func (h *WebHandler) render(c echo.Context, code int, data map[string]any) error {
    c.Response().WriteHeader(code)
    return h.templates.ExecuteTemplate(c.Response().Writer, "base", data)
}
