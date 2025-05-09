package middleware

import (
	"fmt"

	"slices"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/models"
)

func JWT(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired JWT",
			})
		},
	})
}

func RolesMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{
				Error:   "forbidden",
				Message: "Role not found in JWT",
			})
		}

		// for _, allowedRole := range allowedRoles {
		// 	if role == allowedRole {
		// 		return c.Next()
		// 	}
		// }
		if slices.Contains(allowedRoles, role) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{
			Error:   "forbidden",
			Message: fmt.Sprintf("Access restricted to roles: %v", allowedRoles),
		})
	}
}
