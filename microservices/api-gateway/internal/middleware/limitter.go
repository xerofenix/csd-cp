package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/models"
)

func Limiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: "Too many requests, please try again later",
			})
		},
	})
}
