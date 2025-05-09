package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/xerofenix/csd-career/user-service/internal/models"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}
	return c.Status(code).JSON(models.ErrorResponse{
		Error:   "error",
		Message: message,
	})
}
