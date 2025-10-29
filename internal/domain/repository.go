package domain

// FileRepository defines the interface for file storage operations
type FileRepository interface {
	SaveFile(file *UploadedFile) error
	GetFile(id string) (*UploadedFile, error)
	GetFilesByFolder(folderID string) ([]*UploadedFile, error)
	GetAllFiles() ([]*UploadedFile, error)
	DeleteFile(id string) error
}

// FolderRepository defines the interface for folder operations
type FolderRepository interface {
	CreateFolder(folder *Folder) error
	GetFolder(id string) (*Folder, error)
	GetAllFolders() ([]*Folder, error)
	UpdateFolderFileCount(folderID string, count int) error
}

// AdminRepository defines the interface for admin operations
type AdminRepository interface {
	GetAdminByUsername(username string) (*Admin, error)
}
