package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort                string
	DatabaseURL               string
	JWTSecret                 string
	GoogleClientID            string
	GoogleClientSecret        string
	GoogleRedirectURL         string
	SMTPHost                  string
	SMTPPort                  string
	SMTPUser                  string
	SMTPPass                  string
	SMTPFrom                  string
	FrontendURL               string
	VAPIDPublicKey            string
	VAPIDPrivateKey           string
	VAPIDContact              string
	CloudinaryURL             string
	CleanupNotificationsDays  string
	CleanupAppointmentsMonths string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ No .env file found or error loading it, using system environment variables")
	}

	return &Config{
		ServerPort:                getEnv("SERVER_PORT", "8080"),
		DatabaseURL:               getEnv("DATABASE_URL", "postgres://postgres:Nitin%402321@localhost:5432/saloon?sslmode=disable"),
		JWTSecret:                 getEnv("JWT_SECRET", "super-secret-jwt-key-change-in-production"),
		GoogleClientID:            getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:        getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:         getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
		SMTPHost:                  getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:                  getEnv("SMTP_PORT", "587"),
		SMTPUser:                  getEnv("SMTP_USER", ""),
		SMTPPass:                  getEnv("SMTP_PASS", ""),
		SMTPFrom:                  getEnv("SMTP_FROM", "noreply@saloon.com"),
		FrontendURL:               getEnv("FRONTEND_URL", "http://localhost:5173"),
		VAPIDPublicKey:            getEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey:           getEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDContact:              getEnv("VAPID_CONTACT", "mailto:admin@saloon.com"),
		CloudinaryURL:             getEnv("CLOUDINARY_URL", ""),
		CleanupNotificationsDays:  getEnv("CLEANUP_NOTIFICATIONS_DAYS", "30"),
		CleanupAppointmentsMonths: getEnv("CLEANUP_APPOINTMENTS_MONTHS", "6"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
