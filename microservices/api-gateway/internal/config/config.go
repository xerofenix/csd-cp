package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	UserServiceURL         string
	JobServiceURL          string
	AnnouncementServiceURL string
	DashboardServiceURL    string
	NotificationServiceURL string
	JWTSecret              string
	AllowedOrigins         string
	Port                   string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("API_GATEWAY")
	v.AutomaticEnv()

	// Defaults
	v.SetDefault("USER_SERVICE_URL", "http://user-service:8081")
	v.SetDefault("JOB_SERVICE_URL", "http://job-service:8082")
	v.SetDefault("ANNOUNCEMENT_SERVICE_URL", "http://announcement-service:8083")
	v.SetDefault("DASHBOARD_SERVICE_URL", "http://dashboard-service:8084")
	v.SetDefault("NOTIFICATION_SERVICE_URL", "http://notification-service:8085")
	v.SetDefault("JWT_SECRET", "your-secret-key")
	v.SetDefault("ALLOWED_ORIGINS", "http://localhost:3000")
	v.SetDefault("PORT", "8080")

	cfg := &Config{
		UserServiceURL:         v.GetString("USER_SERVICE_URL"),
		JobServiceURL:          v.GetString("JOB_SERVICE_URL"),
		AnnouncementServiceURL: v.GetString("ANNOUNCEMENT_SERVICE_URL"),
		DashboardServiceURL:    v.GetString("DASHBOARD_SERVICE_URL"),
		NotificationServiceURL: v.GetString("NOTIFICATION_SERVICE_URL"),
		JWTSecret:              v.GetString("JWT_SECRET"),
		AllowedOrigins:         v.GetString("ALLOWED_ORIGINS"),
		Port:                   v.GetString("PORT"),
	}

	return cfg, nil
}
