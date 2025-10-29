package handler

import (
	"encoding/json"
	"fileprintapp/internal/usecase"
	ws "fileprintapp/internal/websocket"
	"net/http"
)

// FolderHandler handles folder-related endpoints
type FolderHandler struct {
	folderService *usecase.FolderService
	hub           *ws.Hub
}

// NewFolderHandler creates a new folder handler
func NewFolderHandler(folderService *usecase.FolderService, hub *ws.Hub) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
		hub:           hub,
	}
}

// CreateFolder creates a new folder
func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Folder name is required", http.StatusBadRequest)
		return
	}

	folder, err := h.folderService.CreateFolder(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast to WebSocket clients
	h.hub.BroadcastMessage("folder_created", folder)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

// GetAllFolders retrieves all folders
func (h *FolderHandler) GetAllFolders(w http.ResponseWriter, r *http.Request) {
	folders, err := h.folderService.GetAllFolders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folders)
}
