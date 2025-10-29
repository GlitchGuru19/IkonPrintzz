package postgres

import (
	"database/sql"
	"fileprintapp/internal/domain"
	"time"
)

// FileRepository implements domain.FileRepository using PostgreSQL (Neon)
// This provides persistent storage for uploaded files metadata
type FileRepository struct {
	db *sql.DB // PostgreSQL database connection
}

// NewFileRepository creates a new PostgreSQL-backed file repository
// Parameters:
//   - db: Active database connection to Neon PostgreSQL
// Returns:
//   - Configured FileRepository ready for use
func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

// SaveFile persists file metadata to PostgreSQL database
// This stores information about uploaded files for admin viewing
// Parameters:
//   - file: File entity containing all upload information
// Returns:
//   - error: nil on success, error if database operation fails
func (r *FileRepository) SaveFile(file *domain.UploadedFile) error {
	// SQL query to insert file record into database
	// Uses COALESCE to handle NULL values safely
	query := `
		INSERT INTO uploaded_files (
			id, folder_id, folder_name, file_name, 
			file_size, file_type, file_path, uploaded_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// Set upload timestamp to current time if not already set
	if file.UploadedAt.IsZero() {
		file.UploadedAt = time.Now()
	}

	// Execute INSERT query with file data
	_, err := r.db.Exec(
		query,
		file.ID,
		file.FolderID,
		file.FolderName,
		file.FileName,
		file.FileSize,
		file.FileType,
		file.FilePath,
		file.UploadedAt,
	)

	return err
}

// GetFile retrieves a single file by its unique ID
// Parameters:
//   - id: Unique identifier of the file to retrieve
// Returns:
//   - *domain.UploadedFile: File entity if found
//   - error: sql.ErrNoRows if not found, other errors on query failure
func (r *FileRepository) GetFile(id string) (*domain.UploadedFile, error) {
	query := `
		SELECT id, folder_id, folder_name, file_name, 
		       file_size, file_type, file_path, uploaded_at
		FROM uploaded_files
		WHERE id = $1
	`

	file := &domain.UploadedFile{}
	
	// Scan database row into file struct
	err := r.db.QueryRow(query, id).Scan(
		&file.ID,
		&file.FolderID,
		&file.FolderName,
		&file.FileName,
		&file.FileSize,
		&file.FileType,
		&file.FilePath,
		&file.UploadedAt,
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetFilesByFolder retrieves all files belonging to a specific folder
// Useful for displaying files grouped by user-created folders
// Parameters:
//   - folderID: Unique identifier of the folder
// Returns:
//   - []*domain.UploadedFile: Slice of files in the folder (empty if none)
//   - error: nil on success, error on query failure
func (r *FileRepository) GetFilesByFolder(folderID string) ([]*domain.UploadedFile, error) {
	query := `
		SELECT id, folder_id, folder_name, file_name, 
		       file_size, file_type, file_path, uploaded_at
		FROM uploaded_files
		WHERE folder_id = $1
		ORDER BY uploaded_at DESC
	`

	// Execute query to get all matching rows
	rows, err := r.db.Query(query, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect all files into a slice
	files := make([]*domain.UploadedFile, 0)
	for rows.Next() {
		file := &domain.UploadedFile{}
		err := rows.Scan(
			&file.ID,
			&file.FolderID,
			&file.FolderName,
			&file.FileName,
			&file.FileSize,
			&file.FileType,
			&file.FilePath,
			&file.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

// GetAllFiles retrieves all uploaded files from database
// Ordered by upload time (newest first) for admin dashboard
// Returns:
//   - []*domain.UploadedFile: All files in system
//   - error: nil on success, error on query failure
func (r *FileRepository) GetAllFiles() ([]*domain.UploadedFile, error) {
	query := `
		SELECT id, folder_id, folder_name, file_name, 
		       file_size, file_type, file_path, uploaded_at
		FROM uploaded_files
		ORDER BY uploaded_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]*domain.UploadedFile, 0)
	for rows.Next() {
		file := &domain.UploadedFile{}
		err := rows.Scan(
			&file.ID,
			&file.FolderID,
			&file.FolderName,
			&file.FileName,
			&file.FileSize,
			&file.FileType,
			&file.FilePath,
			&file.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

// DeleteFile removes file metadata from database
// Note: This only deletes database record, not the physical file
// Physical file deletion is handled by the use case layer
// Parameters:
//   - id: Unique identifier of file to delete
// Returns:
//   - error: nil on success, error if file not found or query fails
func (r *FileRepository) DeleteFile(id string) error {
	query := `DELETE FROM uploaded_files WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Return error if file didn't exist
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
