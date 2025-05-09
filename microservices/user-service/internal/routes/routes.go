package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/xerofenix/csd-career/user-service/internal/config"
	"gitlab.com/xerofenix/csd-career/user-service/internal/db"
	"gitlab.com/xerofenix/csd-career/user-service/internal/handlers"
	"gitlab.com/xerofenix/csd-career/user-service/internal/middleware"
	"gitlab.com/xerofenix/csd-career/user-service/internal/storage"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, db *db.DB) {
	userHandler := handlers.NewUserHandler(db, cfg)
	storage, err := storage.New(cfg)
	if err != nil {
		panic(err)
	}
	resumeHandler := handlers.NewResumeHandler(db, storage)

	// Public Routes
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Post("/password-reset", userHandler.PasswordReset)
	app.Post("/verify-email", userHandler.VerifyEmail)

	// Protected Routes
	protected := app.Group("", middleware.JWT(cfg.JWTSecret))
	protected.Get("/users/:id", userHandler.GetProfile)
	protected.Put("/users/:id", userHandler.UpdateProfile)
	protected.Post("/users/resume", resumeHandler.UploadResume)
}
