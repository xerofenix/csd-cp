package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns a Fiber CORS middleware configured for the API Gateway
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // React frontend
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	})
}
