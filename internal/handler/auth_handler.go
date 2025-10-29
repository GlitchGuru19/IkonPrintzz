package handler

import (
	"encoding/json"
	"fileprintapp/internal/usecase"
	"net/http"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *usecase.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles admin login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
