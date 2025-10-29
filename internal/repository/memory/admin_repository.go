package memory

import (
	"errors"
	"fileprintapp/internal/domain"
)

// AdminRepository implements domain.AdminRepository using in-memory storage
type AdminRepository struct {
	admin *domain.Admin
}

// NewAdminRepository creates a new admin repository with credentials
func NewAdminRepository(username, passwordHash string) *AdminRepository {
	return &AdminRepository{
		admin: &domain.Admin{
			Username:     username,
			PasswordHash: passwordHash,
		},
	}
}

// GetAdminByUsername retrieves admin by username
func (r *AdminRepository) GetAdminByUsername(username string) (*domain.Admin, error) {
	if r.admin.Username == username {
		return r.admin, nil
	}
	return nil, errors.New("admin not found")
}
