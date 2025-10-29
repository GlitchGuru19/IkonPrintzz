package postgres

import (
	"database/sql"
	"fileprintapp/internal/domain"
	"time"
)

// FolderRepository implements domain.FolderRepository using PostgreSQL (Neon)
// Manages folder metadata in persistent storage
type FolderRepository struct {
	db *sql.DB // PostgreSQL database connection
}

// NewFolderRepository creates a new PostgreSQL-backed folder repository
// Parameters:
//   - db: Active database connection to Neon PostgreSQL
// Returns:
//   - Configured FolderRepository ready for use
func NewFolderRepository(db *sql.DB) *FolderRepository {
	return &FolderRepository{
		db: db,
	}
}

// CreateFolder persists a new folder to the database
// Folders are used to organize uploaded files by user-defined names
// Parameters:
//   - folder: Folder entity with ID, name, and metadata
// Returns:
//   - error: nil on success, error on duplicate ID or query failure
func (r *FolderRepository) CreateFolder(folder *domain.Folder) error {
	// SQL query to insert folder record
	query := `
		INSERT INTO folders (id, name, created_at, file_count)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`

	// Set creation timestamp if not already set
	if folder.CreatedAt.IsZero() {
		folder.CreatedAt = time.Now()
	}

	// Execute INSERT query
	_, err := r.db.Exec(
		query,
		folder.ID,
		folder.Name,
		folder.CreatedAt,
		folder.FileCount,
	)

	return err
}

// GetFolder retrieves a single folder by its unique ID
// Parameters:
//   - id: Unique identifier of the folder
// Returns:
//   - *domain.Folder: Folder entity if found
//   - error: sql.ErrNoRows if not found, other errors on query failure
func (r *FolderRepository) GetFolder(id string) (*domain.Folder, error) {
	query := `
		SELECT id, name, created_at, file_count
		FROM folders
		WHERE id = $1
	`

	folder := &domain.Folder{}
	
	// Scan database row into folder struct
	err := r.db.QueryRow(query, id).Scan(
		&folder.ID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.FileCount,
	)

	if err != nil {
		return nil, err
	}

	return folder, nil
}

// GetAllFolders retrieves all folders from database
// Ordered by creation time (newest first)
// Returns:
//   - []*domain.Folder: All folders in system
//   - error: nil on success, error on query failure
func (r *FolderRepository) GetAllFolders() ([]*domain.Folder, error) {
	query := `
		SELECT id, name, created_at, file_count
		FROM folders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect all folders into a slice
	folders := make([]*domain.Folder, 0)
	for rows.Next() {
		folder := &domain.Folder{}
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.CreatedAt,
			&folder.FileCount,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, rows.Err()
}

// UpdateFolderFileCount updates the number of files in a folder
// This is called when files are added or removed
// Parameters:
//   - folderID: Unique identifier of the folder
//   - count: New file count (should be >= 0)
// Returns:
//   - error: nil on success, error if folder not found or query fails
func (r *FolderRepository) UpdateFolderFileCount(folderID string, count int) error {
	query := `
		UPDATE folders 
		SET file_count = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, count, folderID)
	if err != nil {
		return err
	}

	// Check if folder exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
