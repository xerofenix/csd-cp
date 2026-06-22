package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/xerofenix/csd-career/user-service/internal/config"
	"gitlab.com/xerofenix/csd-career/user-service/internal/db"
	"gitlab.com/xerofenix/csd-career/user-service/internal/models"
)

type UserHandler struct {
	DB     *db.DB
	Config *config.Config
}

func NewUserHandler(db *db.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{DB: db, Config: cfg}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate role
	validRoles := map[string]bool{"student": true, "tpo": true, "company": true}
	if !validRoles[req.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_role",
			Message: "Role must be student, tpo, or company",
		})
	}

	// Hash password
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
	// 		Error:   "internal_error",
	// 		Message: "Failed to hash password",
	// 	})
	// }

	// Insert user
	var userID int
	err := h.DB.Conn.QueryRow(`
        INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id
    `, req.Email, string(req.Password), req.Role).Scan(&userID)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error:   "user_exists",
			Message: "Email already registered",
		})
	}

	// Insert profile
	detailsJSON, _ := json.Marshal(req.Details)
	_, err = h.DB.Conn.Exec(`
        INSERT INTO profiles (user_id, name, details) VALUES ($1, $2, $3)
    `, userID, req.Name, detailsJSON)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create profile",
		})
	}

	user := models.User{
		ID:    userID,
		Email: req.Email,
		Role:  req.Role,
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	var user models.User
	var hashedPassword string
	err := h.DB.Conn.QueryRow(`
        SELECT id, email, password, role FROM users WHERE email = $1
    `, req.Email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Role)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Invalid email or password",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Database error",
		})
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
	// 		Error:   "invalid_credentials",
	// 		Message: "Invalid email or password",
	// 	})
	// }

	// Generate JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate token",
		})
	}

	return c.JSON(models.LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID",
		})
	}

	var profile models.Profile
	var detailsJSON []byte
	err = h.DB.Conn.QueryRow(`
        SELECT user_id, name, details, updated_at FROM profiles WHERE user_id = $1
    `, id).Scan(&profile.UserID, &profile.Name, &detailsJSON, &profile.UpdatedAt)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "not_found",
			Message: "Profile not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Database error",
		})
	}

	if err := json.Unmarshal(detailsJSON, &profile.Details); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to parse profile details",
		})
	}

	return c.JSON(profile)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID",
		})
	}

	var req models.Profile
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	detailsJSON, _ := json.Marshal(req.Details)
	result, err := h.DB.Conn.Exec(`
        UPDATE profiles SET name = $1, details = $2, updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $3
    `, req.Name, detailsJSON, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update profile",
		})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "not_found",
			Message: "Profile not found",
		})
	}

	req.UserID = id
	req.UpdatedAt = time.Now()
	return c.JSON(req)
}

func (h *UserHandler) PasswordReset(c *fiber.Ctx) error {
	var req models.PasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Placeholder: Generate and send reset code (e.g., via email)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Password reset code sent to %s (placeholder)", req.Email),
	})
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	var req models.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Placeholder: Verify email code
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Email %s verified with code %s (placeholder)", req.Email, req.Code),
	})
}
