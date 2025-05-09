package handlers

import (
	"fmt"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/xerofenix/csd-career/user-service/internal/db"
	"gitlab.com/xerofenix/csd-career/user-service/internal/models"
	"gitlab.com/xerofenix/csd-career/user-service/internal/storage"
)

type ResumeHandler struct {
	DB      *db.DB
	Storage *storage.Storage
}

func NewResumeHandler(db *db.DB, storage *storage.Storage) *ResumeHandler {
	return &ResumeHandler{DB: db, Storage: storage}
}

func (h *ResumeHandler) UploadResume(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["id"].(float64))

	file, err := c.FormFile("resume")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_file",
			Message: "Failed to read resume file",
		})
	}

	// Validate file type (PDF only)
	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_file_type",
			Message: "Only PDF files are allowed",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_file",
			Message: "Failed to open resume file",
		})
	}
	defer src.Close()

	filename := fmt.Sprintf("%d_%s", userID, file.Filename)
	filepath, err := h.Storage.SaveResume(src, filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "storage_error",
			Message: "Failed to save resume",
		})
	}

	var resumeID int
	err = h.DB.Conn.QueryRow(`
        INSERT INTO resumes (user_id, filename, filepath) VALUES ($1, $2, $3) RETURNING id
    `, userID, filename, filepath).Scan(&resumeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to save resume metadata",
		})
	}

	resume := models.Resume{
		ID:         resumeID,
		UserID:     userID,
		Filename:   filename,
		Filepath:   filepath,
		UploadedAt: time.Now(),
	}
	return c.Status(fiber.StatusCreated).JSON(resume)
}
