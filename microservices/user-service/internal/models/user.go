package models

import "time"

// User represents a user in the database
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Excluded from JSON
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Profile represents a user profile
type Profile struct {
	UserID    int            `json:"user_id"`
	Name      string         `json:"name"`
	Details   map[string]any `json:"details"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Resume represents a resume file
type Resume struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Filename   string    `json:"filename"`
	Filepath   string    `json:"filepath"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// RegisterRequest for POST /register
type RegisterRequest struct {
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Role     string         `json:"role"`
	Name     string         `json:"name"`
	Details  map[string]any `json:"details"`
}

// LoginRequest for POST /login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse for POST /login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// PasswordResetRequest for POST /password-reset
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// VerifyEmailRequest for POST /verify-email
type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
