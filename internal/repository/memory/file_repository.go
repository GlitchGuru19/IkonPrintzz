package memory

import (
	"errors"
	"fileprintapp/internal/domain"
	"sync"
)

// FileRepository implements domain.FileRepository using in-memory storage
type FileRepository struct {
	files map[string]*domain.UploadedFile
	mu    sync.RWMutex
}

// NewFileRepository creates a new in-memory file repository
func NewFileRepository() *FileRepository {
	return &FileRepository{
		files: make(map[string]*domain.UploadedFile),
	}
}

// SaveFile saves a file to memory
func (r *FileRepository) SaveFile(file *domain.UploadedFile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.files[file.ID] = file
	return nil
}

// GetFile retrieves a file by ID
func (r *FileRepository) GetFile(id string) (*domain.UploadedFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	file, exists := r.files[id]
	if !exists {
		return nil, errors.New("file not found")
	}
	return file, nil
}

// GetFilesByFolder retrieves all files in a folder
func (r *FileRepository) GetFilesByFolder(folderID string) ([]*domain.UploadedFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var files []*domain.UploadedFile
	for _, file := range r.files {
		if file.FolderID == folderID {
			files = append(files, file)
		}
	}
	return files, nil
}

// GetAllFiles retrieves all files
func (r *FileRepository) GetAllFiles() ([]*domain.UploadedFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	files := make([]*domain.UploadedFile, 0, len(r.files))
	for _, file := range r.files {
		files = append(files, file)
	}
	return files, nil
}

// DeleteFile deletes a file by ID
func (r *FileRepository) DeleteFile(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.files[id]; !exists {
		return errors.New("file not found")
	}
	delete(r.files, id)
	return nil
}
