package main

import (
	"fileprintapp/internal/config"
	"fileprintapp/internal/handler"
	"fileprintapp/internal/middleware"
	"fileprintapp/internal/repository/memory"
	"fileprintapp/internal/usecase"
	ws "fileprintapp/internal/websocket"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(cfg.StoragePath, 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	// Hash admin password
	passwordHash, err := usecase.HashPassword(cfg.AdminPassword)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Initialize repositories
	fileRepo := memory.NewFileRepository()
	folderRepo := memory.NewFolderRepository()
	adminRepo := memory.NewAdminRepository(cfg.AdminUsername, passwordHash)

	// Initialize services
	fileService := usecase.NewFileService(fileRepo, folderRepo, cfg.StoragePath, cfg.MaxFileSize, cfg.AllowedExtensions)
	folderService := usecase.NewFolderService(folderRepo)
	authService := usecase.NewAuthService(adminRepo, cfg.JWTSecret)

	// Initialize WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	fileHandler := handler.NewFileHandler(fileService, folderService, hub)
	folderHandler := handler.NewFolderHandler(folderService, hub)
	wsHandler := handler.NewWebSocketHandler(hub)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup router
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(middleware.CORS)

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Public routes
	r.HandleFunc("/", serveIndexPage).Methods("GET")
	r.HandleFunc("/admin", serveAdminLoginPage).Methods("GET")
	r.HandleFunc("/admin/dashboard", serveAdminDashboardPage).Methods("GET")

	// API routes - Public
	r.HandleFunc("/api/upload", fileHandler.UploadFile).Methods("POST")
	r.HandleFunc("/api/folders", folderHandler.CreateFolder).Methods("POST")
	r.HandleFunc("/api/admin/login", authHandler.Login).Methods("POST")

	// WebSocket route
	r.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// API routes - Protected (Admin only)
	adminRouter := r.PathPrefix("/api").Subrouter()
	adminRouter.Use(authMiddleware.Authenticate)
	adminRouter.HandleFunc("/files", fileHandler.GetAllFiles).Methods("GET")
	adminRouter.HandleFunc("/files/{id}", fileHandler.DeleteFile).Methods("DELETE")
	adminRouter.HandleFunc("/files/{id}/view", fileHandler.ViewFile).Methods("GET")

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("Server starting on http://%s", addr)
	log.Printf("Admin credentials - Username: %s, Password: %s", cfg.AdminUsername, cfg.AdminPassword)
	
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}

func serveAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/admin-login.html")
}

func serveAdminDashboardPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/admin-dashboard.html")
}
