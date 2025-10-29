package usecase

import (
	"fileprintapp/internal/domain"
	"time"

	"github.com/google/uuid"
)

// FolderService handles folder-related business logic
type FolderService struct {
	folderRepo domain.FolderRepository
}

// NewFolderService creates a new folder service
func NewFolderService(folderRepo domain.FolderRepository) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
	}
}

// CreateFolder creates a new folder
func (s *FolderService) CreateFolder(name string) (*domain.Folder, error) {
	folder := &domain.Folder{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		FileCount: 0,
	}

	if err := s.folderRepo.CreateFolder(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// GetAllFolders retrieves all folders
func (s *FolderService) GetAllFolders() ([]*domain.Folder, error) {
	return s.folderRepo.GetAllFolders()
}

// GetFolder retrieves a folder by ID
func (s *FolderService) GetFolder(id string) (*domain.Folder, error) {
	return s.folderRepo.GetFolder(id)
}
