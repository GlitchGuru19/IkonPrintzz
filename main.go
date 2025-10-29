package main

import (
	"fileprintapp/internal/config"
	"fileprintapp/internal/database"
	"fileprintapp/internal/handler"
	"fileprintapp/internal/middleware"
	"fileprintapp/internal/repository/postgres"
	"fileprintapp/internal/usecase"
	ws "fileprintapp/internal/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

// main is the application entry point
// It initializes all components and starts the HTTP server
// Production-ready with PostgreSQL (Neon), graceful shutdown, and proper error handling
func main() {
	// ============================================
	// STEP 1: Load Configuration
	// ============================================
	// Load all settings from environment variables
	// In production, these come from hosting platform (Railway, Fly.io, etc.)
	log.Println("📋 Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Failed to load config:", err)
	}
	log.Printf("✅ Configuration loaded (Environment: %s)", cfg.Environment)

	// ============================================
	// STEP 2: Connect to Database (Neon PostgreSQL)
	// ============================================
	log.Println("🔌 Connecting to Neon PostgreSQL database...")
	
	// Build database configuration from loaded settings
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	// Establish database connection with connection pooling
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	// Ensure database connection is closed on shutdown
	defer database.Close(db)

	// ============================================
	// STEP 3: Run Database Migrations
	// ============================================
	// Create tables and indexes if they don't exist
	// Safe to run on every startup (uses CREATE IF NOT EXISTS)
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	// ============================================
	// STEP 4: Initialize Admin User
	// ============================================
	log.Println("👤 Initializing admin user...")
	
	// Hash the admin password using bcrypt for security
	// Never store passwords in plain text!
	passwordHash, err := usecase.HashPassword(cfg.AdminPassword)
	if err != nil {
		log.Fatal("❌ Failed to hash admin password:", err)
	}

	// Create admin user in database (if not exists)
	if err := database.InitializeAdmin(db, cfg.AdminUsername, passwordHash); err != nil {
		log.Fatal("❌ Failed to initialize admin:", err)
	}

	// ============================================
	// STEP 5: Create Upload Directory
	// ============================================
	// Ensure the uploads directory exists for file storage
	// Even with cloud storage, local temp storage may be needed
	log.Println("📁 Setting up file storage...")
	if err := os.MkdirAll(cfg.StoragePath, 0755); err != nil {
		log.Fatal("❌ Failed to create uploads directory:", err)
	}

	// ============================================
	// STEP 6: Initialize Repositories (Data Layer)
	// ============================================
	// Repositories handle all database operations
	// Using PostgreSQL for persistent storage (replaces in-memory)
	log.Println("💾 Initializing repositories...")
	fileRepo := postgres.NewFileRepository(db)
	folderRepo := postgres.NewFolderRepository(db)
	adminRepo := postgres.NewAdminRepository(db)

	// ============================================
	// STEP 7: Initialize Services (Business Logic Layer)
	// ============================================
	// Services contain business logic and orchestrate repositories
	log.Println("⚙️  Initializing services...")
	fileService := usecase.NewFileService(
		fileRepo,
		folderRepo,
		cfg.StoragePath,
		cfg.MaxFileSize,
		cfg.AllowedExtensions,
	)
	folderService := usecase.NewFolderService(folderRepo)
	authService := usecase.NewAuthService(adminRepo, cfg.JWTSecret)

	// ============================================
	// STEP 8: Initialize WebSocket Hub
	// ============================================
	// WebSocket hub manages real-time connections to admin dashboard
	// Runs in separate goroutine for concurrent handling
	log.Println("🔌 Starting WebSocket hub...")
	hub := ws.NewHub()
	go hub.Run() // Run in background

	// ============================================
	// STEP 9: Initialize HTTP Handlers
	// ============================================
	// Handlers process HTTP requests and responses
	log.Println("🌐 Initializing HTTP handlers...")
	authHandler := handler.NewAuthHandler(authService)
	fileHandler := handler.NewFileHandler(fileService, folderService, hub)
	folderHandler := handler.NewFolderHandler(folderService, hub)
	wsHandler := handler.NewWebSocketHandler(hub)

	// Initialize middleware for cross-cutting concerns
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// ============================================
	// STEP 10: Setup HTTP Router and Routes
	// ============================================
	log.Println("🛣️  Setting up routes...")
	r := mux.NewRouter()

	// Apply CORS middleware to all routes
	// Allows cross-origin requests (important for hosted frontends)
	r.Use(middleware.CORS)

	// Serve static files (HTML, CSS, JS)
	// These are the user-facing web pages
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))),
	)

	// === PUBLIC ROUTES (No authentication required) ===
	
	// Main upload page for users
	r.HandleFunc("/", serveIndexPage).Methods("GET")
	
	// Admin login page (hidden from main page for security)
	// Access via: https://yourdomain.com/admin
	r.HandleFunc("/admin", serveAdminLoginPage).Methods("GET")
	
	// Admin dashboard page
	r.HandleFunc("/admin/dashboard", serveAdminDashboardPage).Methods("GET")

	// API endpoint for file uploads (public - no auth needed)
	r.HandleFunc("/api/upload", fileHandler.UploadFile).Methods("POST")
	
	// API endpoint for creating folders (public)
	r.HandleFunc("/api/folders", folderHandler.CreateFolder).Methods("POST")
	
	// API endpoint for admin login (returns JWT token)
	r.HandleFunc("/api/admin/login", authHandler.Login).Methods("POST")

	// WebSocket endpoint for real-time updates
	r.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// === PROTECTED ROUTES (Require JWT authentication) ===
	// These routes are only accessible to authenticated admins
	adminRouter := r.PathPrefix("/api").Subrouter()
	adminRouter.Use(authMiddleware.Authenticate) // Require valid JWT token
	
	// Get all uploaded files (admin only)
	adminRouter.HandleFunc("/files", fileHandler.GetAllFiles).Methods("GET")
	
	// Delete a file (admin only)
	adminRouter.HandleFunc("/files/{id}", fileHandler.DeleteFile).Methods("DELETE")
	
	// View/print a file (admin only)
	adminRouter.HandleFunc("/files/{id}/view", fileHandler.ViewFile).Methods("GET")

	// ============================================
	// STEP 11: Setup Graceful Shutdown
	// ============================================
	// Listen for interrupt signals (Ctrl+C, SIGTERM from hosting platform)
	// This ensures clean shutdown (close DB connections, finish requests, etc.)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to handle shutdown
	go func() {
		<-shutdownChan // Wait for shutdown signal
		log.Println("\n🛑 Shutdown signal received, cleaning up...")
		database.Close(db) // Close database connection
		os.Exit(0)
	}()

	// ============================================
	// STEP 12: Start HTTP Server
	// ============================================
	// Build server address from configuration
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	
	// Log startup information
	separator := repeat("=", 60)
	log.Println("\n" + separator)
	log.Println("🚀 File Print Service - PRODUCTION MODE")
	log.Println(separator)
	log.Printf("📍 Server Address: http://%s", addr)
	log.Printf("📁 User Upload Page: http://%s", addr)
	log.Printf("🔐 Admin Login: http://%s/admin", addr)
	log.Println("\n⚠️  ADMIN ACCESS INSTRUCTIONS:")
	log.Println("   When hosted, access admin at: https://yourdomain.com/admin")
	log.Println("   The admin link is NOT shown on the user page for security")
	log.Printf("   Username: %s\n", cfg.AdminUsername)
	log.Println(separator + "\n")
	
	// Start the HTTP server
	// This blocks until server is stopped
	log.Printf("✅ Server is running and ready to accept connections!")
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}

// Helper function to repeat strings (for visual formatting)
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
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
