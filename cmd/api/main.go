package main

import (
	"context"
	"log"

	"fidely-backend/internal/config"
	"fidely-backend/internal/db"
	"fidely-backend/internal/handler"

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

	// Initialize web handler with templates
	webHandler, err := handler.NewWebHandler()
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

	// API Routes will be added here
	// e.g., e.GET("/api/stores", storeHandler.List)

	// Start server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
