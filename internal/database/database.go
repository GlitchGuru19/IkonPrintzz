package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database connection configuration
// All values should come from environment variables for security
type Config struct {
	Host     string // Database host (from Neon)
	Port     string // Database port (usually 5432)
	User     string // Database user
	Password string // Database password (keep secret!)
	DBName   string // Database name
	SSLMode  string // SSL mode (require for Neon)
}

// Connect establishes a connection to PostgreSQL (Neon)
// It configures connection pooling and validates the connection
// Parameters:
//   - cfg: Database configuration with connection details
// Returns:
//   - *sql.DB: Active database connection pool
//   - error: nil on success, error if connection fails
func Connect(cfg Config) (*sql.DB, error) {
	// Build PostgreSQL connection string
	// Format: postgres://user:password@host:port/dbname?sslmode=require
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Open database connection
	// Note: This doesn't actually connect yet, just prepares the driver
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for optimal performance
	// MaxOpenConns: Maximum number of open connections to the database
	// Neon has connection limits, so we keep this moderate
	db.SetMaxOpenConns(25)

	// MaxIdleConns: Maximum number of idle connections in the pool
	// Keeping some idle connections reduces latency for new requests
	db.SetMaxIdleConns(5)

	// ConnMaxLifetime: Maximum time a connection can be reused
	// Important for Neon to handle connection rotation
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection is actually working
	// This is the actual connection attempt
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established successfully")
	return db, nil
}

// RunMigrations executes database schema migrations
// This should be run on application startup to ensure schema is up-to-date
// Parameters:
//   - db: Active database connection
// Returns:
//   - error: nil on success, error if migrations fail
func RunMigrations(db *sql.DB) error {
	log.Println("🔄 Running database migrations...")

	// Migration: Create folders table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS folders (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			file_count INTEGER NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create folders table: %w", err)
	}

	// Migration: Create uploaded_files table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS uploaded_files (
			id VARCHAR(255) PRIMARY KEY,
			folder_id VARCHAR(255) NOT NULL,
			folder_name VARCHAR(255) NOT NULL,
			file_name VARCHAR(500) NOT NULL,
			file_size BIGINT NOT NULL,
			file_type VARCHAR(50) NOT NULL,
			file_path TEXT NOT NULL,
			uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create uploaded_files table: %w", err)
	}

	// Migration: Create admins table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS admins (
			username VARCHAR(100) PRIMARY KEY,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create admins table: %w", err)
	}

	// Migration: Create indexes for better performance
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_folders_created_at ON folders(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_uploaded_files_folder_id ON uploaded_files(folder_id);
		CREATE INDEX IF NOT EXISTS idx_uploaded_files_uploaded_at ON uploaded_files(uploaded_at DESC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// InitializeAdmin creates the default admin user if it doesn't exist
// This ensures there's always at least one admin account
// Parameters:
//   - db: Active database connection
//   - username: Admin username
//   - passwordHash: Bcrypt-hashed password
// Returns:
//   - error: nil on success, error if operation fails
func InitializeAdmin(db *sql.DB, username, passwordHash string) error {
	query := `
		INSERT INTO admins (username, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (username) DO NOTHING
	`

	_, err := db.Exec(query, username, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to initialize admin: %w", err)
	}

	log.Printf("✅ Admin user '%s' initialized", username)
	return nil
}

// Close gracefully closes the database connection
// Should be called on application shutdown
// Parameters:
//   - db: Database connection to close
func Close(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("⚠️  Error closing database: %v", err)
		} else {
			log.Println("✅ Database connection closed")
		}
	}
}
