package tests

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/middleware"
)

func TestRolesMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.JWT("test-secret"))
	app.Use(middleware.RolesMiddleware("student"))
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("OK") })

	// Create a valid student token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "student",
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
