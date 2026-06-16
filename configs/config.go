package configs

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port        string
	DatabaseURL string
	AppEnv      string
}

// Load reads configuration from environment variables with sensible defaults.
// It attempts to load a .env file if present but does not fail if missing,
// since production environments typically inject env vars directly.
func Load() *Config {
	// Best-effort .env load — ignore error for production deployments.
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ainyx?sslmode=disable"),
		AppEnv:      getEnv("APP_ENV", "development"),
	}
}

// getEnv retrieves an environment variable or returns a fallback default.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
