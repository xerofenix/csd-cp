package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	UploadDir   string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("USER_SERVICE")
	v.AutomaticEnv()

	// Defaults
	v.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/career_db?sslmode=disable")
	v.SetDefault("JWT_SECRET", "your-secret-key")
	v.SetDefault("PORT", "8081")
	v.SetDefault("UPLOAD_DIR", "./uploads/resumes")

	cfg := &Config{
		DatabaseURL: v.GetString("DATABASE_URL"),
		JWTSecret:   v.GetString("JWT_SECRET"),
		Port:        v.GetString("PORT"),
		UploadDir:   v.GetString("UPLOAD_DIR"),
	}

	return cfg, nil
}
