package handler

import (
	"encoding/json"
	"fileprintapp/internal/usecase"
	ws "fileprintapp/internal/websocket"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// FileHandler handles file-related endpoints
type FileHandler struct {
	fileService   *usecase.FileService
	folderService *usecase.FolderService
	hub           *ws.Hub
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileService *usecase.FileService, folderService *usecase.FolderService, hub *ws.Hub) *FileHandler {
	return &FileHandler{
		fileService:   fileService,
		folderService: folderService,
		hub:           hub,
	}
}

// UploadFile handles file upload
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	folderID := r.FormValue("folder_id")
	folderName := r.FormValue("folder_name")

	if folderID == "" || folderName == "" {
		http.Error(w, "Folder ID and name are required", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	file.Close()

	// Upload file
	uploadedFile, err := h.fileService.UploadFile(handler, folderID, folderName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast to WebSocket clients
	h.hub.BroadcastMessage("new_file", uploadedFile)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadedFile)
}

// GetAllFiles retrieves all files
func (h *FileHandler) GetAllFiles(w http.ResponseWriter, r *http.Request) {
	files, err := h.fileService.GetAllFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// DeleteFile deletes a file
func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["id"]

	if err := h.fileService.DeleteFile(fileID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast to WebSocket clients
	h.hub.BroadcastMessage("file_deleted", map[string]string{"id": fileID})

	w.WriteHeader(http.StatusNoContent)
}

// ViewFile serves a file for viewing/printing
func (h *FileHandler) ViewFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["id"]

	file, err := h.fileService.GetFile(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate content type
	ext := filepath.Ext(file.FileName)
	switch ext {
	case ".pdf":
		w.Header().Set("Content-Type", "application/pdf")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Set headers to allow inline viewing (not download)
	w.Header().Set("Content-Disposition", "inline; filename=\""+file.FileName+"\"")

	http.ServeFile(w, r, file.FilePath)
}
