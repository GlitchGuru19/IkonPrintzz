package postgres

import (
	"database/sql"
	"fileprintapp/internal/domain"
)

// AdminRepository implements domain.AdminRepository using PostgreSQL (Neon)
// Manages admin user authentication data in persistent storage
type AdminRepository struct {
	db *sql.DB // PostgreSQL database connection
}

// NewAdminRepository creates a new PostgreSQL-backed admin repository
// Parameters:
//   - db: Active database connection to Neon PostgreSQL
// Returns:
//   - Configured AdminRepository ready for use
func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

// GetAdminByUsername retrieves admin credentials by username
// Used during login to verify credentials
// Parameters:
//   - username: Admin username to look up
// Returns:
//   - *domain.Admin: Admin entity with hashed password
//   - error: sql.ErrNoRows if not found, other errors on query failure
func (r *AdminRepository) GetAdminByUsername(username string) (*domain.Admin, error) {
	query := `
		SELECT username, password_hash
		FROM admins
		WHERE username = $1
	`

	admin := &domain.Admin{}
	
	// Scan database row into admin struct
	err := r.db.QueryRow(query, username).Scan(
		&admin.Username,
		&admin.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return admin, nil
}

// CreateAdmin creates a new admin user in database
// Used during initial setup or when adding new admins
// Parameters:
//   - admin: Admin entity with username and hashed password
// Returns:
//   - error: nil on success, error on duplicate username or query failure
func (r *AdminRepository) CreateAdmin(admin *domain.Admin) error {
	query := `
		INSERT INTO admins (username, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (username) DO NOTHING
	`

	_, err := r.db.Exec(query, admin.Username, admin.PasswordHash)
	return err
}
