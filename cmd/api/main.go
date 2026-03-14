package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"fidely-backend/internal/auth"
	"fidely-backend/internal/config"
	"fidely-backend/internal/db"
	"fidely-backend/internal/handler"
	"fidely-backend/internal/repository"
	"fidely-backend/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	authRepository := repository.NewAdminAuthRepository(pool)
	passwordManager := auth.NewDefaultPasswordManager()
	sessionManager := auth.NewSessionManager(cfg.AuthTokenHashPepper)
	authService, err := service.NewAdminAuthService(authRepository, passwordManager, sessionManager, cfg.AuthSessionTTL)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}
	authMiddleware := handler.NewAuthMiddleware(cfg, authService)
	loginRateLimiter, err := handler.NewRedisLoginRateLimiter(cfg.RedisURL, cfg.LoginRateLimitMax, cfg.LoginRateLimitWindow)
	if err != nil {
		log.Fatalf("Failed to initialize login rate limiter: %v", err)
	}
	defer func() {
		if closeErr := loginRateLimiter.Close(); closeErr != nil {
			log.Printf("Failed to close login rate limiter: %v", closeErr)
		}
	}()

	// Initialize web handler with templates
	webHandler, err := handler.NewWebHandler(cfg, authService)
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:_csrf",
		CookieHTTPOnly: true,
		CookieSecure:   cfg.AuthCookieSecure,
		CookieSameSite: toEchoSameSiteMode(cfg.AuthCookieSameSite),
		CookiePath:     "/",
	}))

	// Static files (CSS, JS, images)
	e.Static("/static", "web/static")

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()

		postgresErr := pool.Ping(ctx)
		redisErr := loginRateLimiter.Ping(ctx)

		dependencies := map[string]string{
			"postgres": "ok",
			"redis":    "ok",
		}
		if postgresErr != nil {
			dependencies["postgres"] = "down"
		}
		if redisErr != nil {
			dependencies["redis"] = "down"
		}

		status := "ok"
		httpStatus := http.StatusOK
		if postgresErr != nil || redisErr != nil {
			status = "degraded"
			httpStatus = http.StatusServiceUnavailable
		}

		return c.JSON(httpStatus, map[string]any{
			"status":       status,
			"dependencies": dependencies,
		})
	})

	// Web Routes (Admin Interface)
	e.GET("/", webHandler.LoginPage)
	e.GET("/login", webHandler.LoginPage)
	e.POST("/auth/login", webHandler.HandleLogin, loginRateLimiter.Middleware())

	authenticated := e.Group("/auth")
	authenticated.Use(authMiddleware.RequireAuthenticatedAdmin())
	authenticated.POST("/logout", webHandler.HandleLogout)
	authenticated.GET("/me", webHandler.CurrentAdmin)

	platformOnly := e.Group("/admin/platform")
	platformOnly.Use(authMiddleware.RequireAuthenticatedAdmin(), authMiddleware.RequireAdminTypes(auth.AdminTypeFidelyAdmin))
	platformOnly.GET("/status", webHandler.PlatformOnlyStatus)

	// API Routes will be added here
	// e.g., e.GET("/api/stores", storeHandler.List)

	// Start server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func toEchoSameSiteMode(value string) http.SameSite {
	switch value {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
