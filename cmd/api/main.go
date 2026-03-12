package main

import (
	"context"
	"log"

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

	// Static files (CSS, JS, images)
	e.Static("/static", "web/static")

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Web Routes (Admin Interface)
	e.GET("/", webHandler.LoginPage)
	e.GET("/login", webHandler.LoginPage)
	e.POST("/auth/login", webHandler.HandleLogin)

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
