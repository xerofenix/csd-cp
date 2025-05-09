package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/config"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/middleware"
	"gitlab.com/xerofenix/csd-career/api-gateway/internal/proxy"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Public Routes
	public := app.Group("/api")
	public.Post("/login", proxy.Proxy(cfg.UserServiceURL, "/login"))
	public.Post("/register", proxy.Proxy(cfg.UserServiceURL, "/register"))

	// Protected Routes
	protected := app.Group("/api")
	protected.Use(middleware.JWT(cfg.JWTSecret))

	// User Service Routes
	protected.Get("/users/:id", proxy.Proxy(cfg.UserServiceURL, "/users/:id"))
	protected.Put("/users/:id", proxy.Proxy(cfg.UserServiceURL, "/users/:id"))
	protected.Post("/users/resume", middleware.RolesMiddleware("student"), proxy.Proxy(cfg.UserServiceURL, "/users/resume"))

	// Job Service Routes
	protected.Get("/jobs", proxy.Proxy(cfg.JobServiceURL, "/jobs"))
	protected.Get("/jobs/:id", proxy.Proxy(cfg.JobServiceURL, "/jobs/:id"))
	protected.Post("/jobs", middleware.RolesMiddleware("company"), proxy.Proxy(cfg.JobServiceURL, "/jobs"))
	protected.Post("/jobs/:id/apply", middleware.RolesMiddleware("student"), proxy.Proxy(cfg.JobServiceURL, "/jobs/:id/apply"))
	protected.Get("/jobs/:id/applicants", middleware.RolesMiddleware("company"), proxy.Proxy(cfg.JobServiceURL, "/jobs/:id/applicants"))

	// Announcement Service Routes
	protected.Get("/announcements", proxy.Proxy(cfg.AnnouncementServiceURL, "/announcements"))
	protected.Post("/announcements", middleware.RolesMiddleware("company", "tpo"), proxy.Proxy(cfg.AnnouncementServiceURL, "/announcements"))

	// Dashboard Service Routes
	protected.Get("/stats/companies", middleware.RolesMiddleware("tpo"), proxy.Proxy(cfg.DashboardServiceURL, "/stats/companies"))
	protected.Get("/stats/jobs", middleware.RolesMiddleware("tpo"), proxy.Proxy(cfg.DashboardServiceURL, "/stats/jobs"))
	protected.Get("/stats/applications", middleware.RolesMiddleware("tpo"), proxy.Proxy(cfg.DashboardServiceURL, "/stats/applications"))

	// Notification Service Routes
	protected.Get("/notifications", proxy.Proxy(cfg.NotificationServiceURL, "/notifications"))
}
