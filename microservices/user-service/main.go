package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/xerofenix/csd-career/user-service/internal/config"
	"gitlab.com/xerofenix/csd-career/user-service/internal/db"
	"gitlab.com/xerofenix/csd-career/user-service/internal/handlers"
	"gitlab.com/xerofenix/csd-career/user-service/internal/middleware"
	"gitlab.com/xerofenix/csd-career/user-service/internal/routes"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(middleware.Logger())

	// Setup routes
	routes.SetupRoutes(app, cfg, database)

	// Health check
	app.Get("/health", handlers.HealthCheck())

	// Start server
	port := cfg.Port
	log.Fatal(app.Listen(":" + port))
}
