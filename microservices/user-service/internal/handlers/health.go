package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func HealthCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"message": "User Service is running",
		})
	}
}
