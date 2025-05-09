package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/config"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/health"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/middleware"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/routes"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())
	app.Use(middleware.Limiter())
	middleware.SetupPrometheus(app)

	// Routes
	routes.SetupRoutes(app, cfg)

	// Health check
	app.Get("/health", health.Check(cfg))

	// Start server
	port := cfg.Port
	log.Fatal(app.Listen(":" + port))
}
