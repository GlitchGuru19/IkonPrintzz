package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
// All values are loaded from environment variables for security
// In production, set these in your hosting platform's environment settings
type Config struct {
	// Server configuration
	Port        string // Port to run server on (e.g., "8080")
	Host        string // Host to bind to ("0.0.0.0" for production, "localhost" for dev)
	Environment string // "development" or "production"

	// Admin authentication
	AdminUsername string // Admin username for dashboard access
	AdminPassword string // Admin password (will be hashed with bcrypt)
	JWTSecret     string // Secret key for signing JWT tokens (keep secure!)

	// File upload settings
	MaxFileSize       int64    // Maximum file size in bytes
	AllowedExtensions []string // Allowed file extensions (e.g., ["pdf", "jpg"])

	// Storage configuration
	StorageType string // "local" or "cloud"
	StoragePath string // Path for local storage or cloud config

	// Database configuration (Neon PostgreSQL)
	DBHost     string // Database host from Neon
	DBPort     string // Database port (usually 5432)
	DBName     string // Database name
	DBUser     string // Database username
	DBPassword string // Database password
	DBSSLMode  string // SSL mode ("require" for Neon)
}

// LoadConfig loads configuration from environment variables
// It first attempts to load from .env file, then falls back to system env vars
// This is safe for production as hosting platforms set environment variables
// Returns:
//   - *Config: Fully populated configuration object
//   - error: Currently always nil, but kept for future validation
func LoadConfig() (*Config, error) {
	// Attempt to load .env file
	// In production (Railway, Fly.io, etc.), this file won't exist
	// and env vars will be loaded from the platform instead
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  No .env file found, using system environment variables")
	}

	// Parse max file size from string to int64
	// Default: 10MB (10485760 bytes)
	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "10485760"), 10, 64)
	
	// Parse allowed extensions from comma-separated string
	// Default: common image and document formats
	extensions := strings.Split(getEnv("ALLOWED_EXTENSIONS", "jpg,jpeg,png,pdf,gif"), ",")

	// Build and return configuration object
	return &Config{
		// Server settings
		Port:        getEnv("PORT", "8080"),
		Host:        getEnv("HOST", "0.0.0.0"), // 0.0.0.0 allows external connections
		Environment: getEnv("ENVIRONMENT", "production"),

		// Admin authentication
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "changeme123"),
		JWTSecret:     getEnv("JWT_SECRET", "change-this-secret-key-in-production"),

		// File upload configuration
		MaxFileSize:       maxFileSize,
		AllowedExtensions: extensions,

		// Storage settings
		StorageType: getEnv("STORAGE_TYPE", "local"),
		StoragePath: getEnv("STORAGE_PATH", "./uploads"),

		// Database configuration (Neon PostgreSQL)
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "neondb"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSL_MODE", "require"), // Always require for Neon
	}, nil
}

// getEnv retrieves an environment variable or returns a default value
// This is a helper function to safely get env vars with fallbacks
// Parameters:
//   - key: Environment variable name
//   - defaultValue: Value to return if env var is not set
// Returns:
//   - string: Environment variable value or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment checks if we're running in development mode
// Useful for conditional logging, debugging, etc.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction checks if we're running in production mode
// Use this to enable production-only features
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
