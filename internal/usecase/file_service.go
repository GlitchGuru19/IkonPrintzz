package usecase

import (
	"errors"
	"fileprintapp/internal/domain"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// FileService handles file-related business logic
type FileService struct {
	fileRepo      domain.FileRepository
	folderRepo    domain.FolderRepository
	uploadPath    string
	maxFileSize   int64
	allowedExtensions []string
}

// NewFileService creates a new file service
func NewFileService(fileRepo domain.FileRepository, folderRepo domain.FolderRepository, uploadPath string, maxFileSize int64, allowedExtensions []string) *FileService {
	return &FileService{
		fileRepo:      fileRepo,
		folderRepo:    folderRepo,
		uploadPath:    uploadPath,
		maxFileSize:   maxFileSize,
		allowedExtensions: allowedExtensions,
	}
}

// UploadFile handles file upload logic
func (s *FileService) UploadFile(fileHeader *multipart.FileHeader, folderID, folderName string) (*domain.UploadedFile, error) {
	// Validate file size
	if fileHeader.Size > s.maxFileSize {
		return nil, errors.New("file size exceeds maximum allowed size")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	ext = strings.TrimPrefix(ext, ".")
	if !s.isAllowedExtension(ext) {
		return nil, errors.New("file type not allowed")
	}

	// Open uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create folder directory if it doesn't exist
	folderPath := filepath.Join(s.uploadPath, folderID)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, err
	}

	// Generate unique file ID and path
	fileID := uuid.New().String()
	fileName := fileHeader.Filename
	filePath := filepath.Join(folderPath, fileID+filepath.Ext(fileName))

	// Save file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	// Create file entity
	uploadedFile := &domain.UploadedFile{
		ID:         fileID,
		FolderID:   folderID,
		FolderName: folderName,
		FileName:   fileName,
		FileSize:   fileHeader.Size,
		FileType:   ext,
		FilePath:   filePath,
	}

	// Save to repository
	if err := s.fileRepo.SaveFile(uploadedFile); err != nil {
		os.Remove(filePath) // Clean up file on error
		return nil, err
	}

	// Update folder file count
	files, _ := s.fileRepo.GetFilesByFolder(folderID)
	s.folderRepo.UpdateFolderFileCount(folderID, len(files))

	return uploadedFile, nil
}

// GetAllFiles retrieves all uploaded files
func (s *FileService) GetAllFiles() ([]*domain.UploadedFile, error) {
	return s.fileRepo.GetAllFiles()
}

// GetFilesByFolder retrieves files by folder ID
func (s *FileService) GetFilesByFolder(folderID string) ([]*domain.UploadedFile, error) {
	return s.fileRepo.GetFilesByFolder(folderID)
}

// DeleteFile deletes a file
func (s *FileService) DeleteFile(fileID string) error {
	file, err := s.fileRepo.GetFile(fileID)
	if err != nil {
		return err
	}

	// Delete physical file
	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Delete from repository
	if err := s.fileRepo.DeleteFile(fileID); err != nil {
		return err
	}

	// Update folder file count
	files, _ := s.fileRepo.GetFilesByFolder(file.FolderID)
	s.folderRepo.UpdateFolderFileCount(file.FolderID, len(files))

	return nil
}

// GetFile retrieves a file by ID
func (s *FileService) GetFile(fileID string) (*domain.UploadedFile, error) {
	return s.fileRepo.GetFile(fileID)
}

func (s *FileService) isAllowedExtension(ext string) bool {
	for _, allowed := range s.allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}
