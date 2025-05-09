package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/models"
)

func Proxy(serviceURL, path string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetPath := strings.Replace(path, ":id", c.Params("id"), -1)
		targetURL := fmt.Sprintf("%s%s", serviceURL, targetPath)

		// Convert body to io.Reader
		var bodyReader io.Reader
		if body := c.Request().Body(); len(body) > 0 {
			bodyReader = bytes.NewReader(body)
		} else {
			bodyReader = nil
		}

		req, err := http.NewRequest(c.Method(), targetURL, bodyReader)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "proxy_error",
				Message: "Failed to create request",
			})
		}

		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		client := &http.Client{Timeout: 10 * time.Second} // Configurable timeout
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(models.ErrorResponse{
				Error:   "service_unavailable",
				Message: "Failed to reach service",
			})
		}
		defer resp.Body.Close()

		var body any
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "invalid_response",
				Message: "Failed to parse service response",
			})
		}

		return c.Status(resp.StatusCode).JSON(body)
	}
}
