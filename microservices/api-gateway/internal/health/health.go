package health

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/config"
)

func Check(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		services := map[string]string{
			"user":         cfg.UserServiceURL,
			"job":          cfg.JobServiceURL,
			"announcement": cfg.AnnouncementServiceURL,
			"dashboard":    cfg.DashboardServiceURL,
			"notification": cfg.NotificationServiceURL,
		}

		status := make(map[string]string)
		client := &http.Client{Timeout: 2 * time.Second}

		for name, url := range services {
			resp, err := client.Get(url + "/health") // Assumes each service has a /health endpoint
			if err != nil || resp.StatusCode != http.StatusOK {
				status[name] = "unhealthy"
			} else {
				status[name] = "healthy"
			}
			if resp != nil {
				resp.Body.Close()
			}
		}

		overallStatus := http.StatusOK
		for _, s := range status {
			if s == "unhealthy" {
				overallStatus = http.StatusServiceUnavailable
				break
			}
		}

		return c.Status(overallStatus).JSON(map[string]any{
			"status":  status,
			"message": "API Gateway health check",
		})
	}
}
