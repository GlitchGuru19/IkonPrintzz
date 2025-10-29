package domain

import "time"

// UploadedFile represents a file uploaded by a user
type UploadedFile struct {
	ID         string    `json:"id"`
	FolderID   string    `json:"folder_id"`
	FolderName string    `json:"folder_name"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `json:"file_type"`
	FilePath   string    `json:"file_path"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Folder represents a collection of files
type Folder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	FileCount int       `json:"file_count"`
}

// Admin represents an admin user
type Admin struct {
	Username     string
	PasswordHash string
}

// WebSocketMessage represents a message sent via WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"` // "new_file", "file_deleted", "folder_created"
	Payload interface{} `json:"payload"`
}
