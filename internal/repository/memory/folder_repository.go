package memory

import (
	"errors"
	"fileprintapp/internal/domain"
	"sync"
)

// FolderRepository implements domain.FolderRepository using in-memory storage
type FolderRepository struct {
	folders map[string]*domain.Folder
	mu      sync.RWMutex
}

// NewFolderRepository creates a new in-memory folder repository
func NewFolderRepository() *FolderRepository {
	return &FolderRepository{
		folders: make(map[string]*domain.Folder),
	}
}

// CreateFolder creates a new folder
func (r *FolderRepository) CreateFolder(folder *domain.Folder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.folders[folder.ID] = folder
	return nil
}

// GetFolder retrieves a folder by ID
func (r *FolderRepository) GetFolder(id string) (*domain.Folder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	folder, exists := r.folders[id]
	if !exists {
		return nil, errors.New("folder not found")
	}
	return folder, nil
}

// GetAllFolders retrieves all folders
func (r *FolderRepository) GetAllFolders() ([]*domain.Folder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	folders := make([]*domain.Folder, 0, len(r.folders))
	for _, folder := range r.folders {
		folders = append(folders, folder)
	}
	return folders, nil
}

// UpdateFolderFileCount updates the file count for a folder
func (r *FolderRepository) UpdateFolderFileCount(folderID string, count int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	folder, exists := r.folders[folderID]
	if !exists {
		return errors.New("folder not found")
	}
	folder.FileCount = count
	return nil
}
