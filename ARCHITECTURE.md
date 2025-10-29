# Architecture Documentation

## 📐 Clean Architecture Overview

This application follows Clean Architecture principles with clear separation of concerns.

```
┌─────────────────────────────────────────┐
│         Presentation Layer              │
│  (Handlers, Middleware, WebSocket)      │
├─────────────────────────────────────────┤
│         Use Case Layer                  │
│  (Business Logic, Services)             │
├─────────────────────────────────────────┤
│         Repository Layer                │
│  (Data Access, Storage)                 │
├─────────────────────────────────────────┤
│         Domain Layer                    │
│  (Entities, Interfaces)                 │
└─────────────────────────────────────────┘
```

### Dependency Rule

Dependencies point **inward only**:
- **Domain** has NO dependencies
- **Repository** depends on Domain
- **Use Case** depends on Domain (interfaces)
- **Handler** depends on Use Case and Domain

## 🗂️ Layer Details

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Core business entities and interfaces

**Files**:
- `entities.go` - Business entities (UploadedFile, Folder, Admin)
- `repository.go` - Repository interfaces

**Key Points**:
- No external dependencies
- Defines what the system does (entities)
- Defines how to access data (interfaces)
- Other layers implement these interfaces

**Example**:
```go
type UploadedFile struct {
    ID         string
    FolderID   string
    FileName   string
    // ...
}

type FileRepository interface {
    SaveFile(file *UploadedFile) error
    GetFile(id string) (*UploadedFile, error)
    // ...
}
```

### 2. Repository Layer (`internal/repository/`)

**Purpose**: Data storage implementations

**Structure**:
```
repository/
└── memory/
    ├── file_repository.go
    ├── folder_repository.go
    └── admin_repository.go
```

**Key Points**:
- Implements domain repository interfaces
- Currently uses in-memory storage
- Easy to swap for database (PostgreSQL, MongoDB, etc.)
- Thread-safe with sync.RWMutex

**Example**:
```go
type FileRepository struct {
    files map[string]*domain.UploadedFile
    mu    sync.RWMutex
}

func (r *FileRepository) SaveFile(file *domain.UploadedFile) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.files[file.ID] = file
    return nil
}
```

### 3. Use Case Layer (`internal/usecase/`)

**Purpose**: Business logic and application services

**Files**:
- `file_service.go` - File management logic
- `folder_service.go` - Folder management logic
- `auth_service.go` - Authentication logic

**Key Points**:
- Orchestrates data flow
- Implements business rules
- Uses repository interfaces (dependency injection)
- No HTTP or WebSocket knowledge

**Example**:
```go
type FileService struct {
    fileRepo      domain.FileRepository
    folderRepo    domain.FolderRepository
    uploadPath    string
    maxFileSize   int64
    // ...
}

func (s *FileService) UploadFile(fileHeader *multipart.FileHeader, folderID, folderName string) (*domain.UploadedFile, error) {
    // Validate file
    // Save to disk
    // Save to repository
    // Update folder count
    return uploadedFile, nil
}
```

### 4. Handler Layer (`internal/handler/`)

**Purpose**: HTTP request handling

**Files**:
- `file_handler.go` - File endpoints
- `folder_handler.go` - Folder endpoints
- `auth_handler.go` - Authentication endpoints
- `websocket_handler.go` - WebSocket connections

**Key Points**:
- Converts HTTP requests to use case calls
- Converts responses to JSON
- Broadcasts WebSocket messages
- No business logic

**Example**:
```go
type FileHandler struct {
    fileService   *usecase.FileService
    folderService *usecase.FolderService
    hub           *ws.Hub
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Call service
    // Broadcast update
    // Return response
}
```

### 5. Middleware Layer (`internal/middleware/`)

**Purpose**: Cross-cutting concerns

**Files**:
- `auth_middleware.go` - JWT authentication
- `cors_middleware.go` - CORS headers

**Key Points**:
- Wraps HTTP handlers
- Validates JWT tokens
- Adds context values
- CORS configuration

### 6. WebSocket Layer (`internal/websocket/`)

**Purpose**: Real-time communication

**Files**:
- `hub.go` - Connection manager
- `client.go` - Client connection handler

**Key Points**:
- Manages active connections
- Broadcasts to all clients
- Handles connect/disconnect
- Thread-safe with mutexes

**Flow**:
```
Admin connects → Hub registers client
File uploaded → Handler broadcasts → Hub sends to all clients
Admin disconnects → Hub unregisters client
```

### 7. Configuration Layer (`internal/config/`)

**Purpose**: Application configuration

**Files**:
- `config.go` - Config loading from .env

**Key Points**:
- Loads environment variables
- Provides defaults
- Type conversions
- Single source of truth

## 🔄 Data Flow Examples

### File Upload Flow

```
1. User (Browser)
   ↓ POST /api/upload
2. FileHandler.UploadFile
   ↓ calls
3. FileService.UploadFile
   ↓ validates & saves to disk
   ↓ calls
4. FileRepository.SaveFile
   ↓ stores in memory
   ↓ returns
5. FileHandler.UploadFile
   ↓ broadcasts
6. WebSocket Hub
   ↓ sends to all
7. Admin Dashboard (Browser)
```

### Admin Authentication Flow

```
1. Admin (Browser)
   ↓ POST /api/admin/login
2. AuthHandler.Login
   ↓ calls
3. AuthService.Login
   ↓ calls
4. AdminRepository.GetAdminByUsername
   ↓ returns admin
   ↓ verifies password
   ↓ generates JWT
5. AuthHandler.Login
   ↓ returns token
6. Admin (Browser)
   ↓ stores token
   ↓ uses for protected routes
```

### Protected Route Access

```
1. Admin (Browser)
   ↓ GET /api/files (with Authorization header)
2. AuthMiddleware.Authenticate
   ↓ validates token
   ↓ adds username to context
3. FileHandler.GetAllFiles
   ↓ calls
4. FileService.GetAllFiles
   ↓ calls
5. FileRepository.GetAllFiles
   ↓ returns files
6. Admin (Browser)
```

## 🎯 Design Patterns Used

### 1. Repository Pattern
- Abstracts data access
- Easy to swap implementations
- Testable with mock repositories

### 2. Dependency Injection
- Services receive dependencies via constructors
- Loose coupling
- Easy to test

### 3. Interface Segregation
- Small, focused interfaces
- Domain defines interfaces
- Implementations in other layers

### 4. Pub/Sub Pattern (WebSocket)
- Hub manages subscriptions
- Broadcasts to all subscribers
- Decoupled communication

## 🧪 Testing Strategy

### Unit Tests (Recommended)

**Domain Layer**:
- Test entity validation
- No mocks needed

**Use Case Layer**:
- Mock repositories
- Test business logic
- Test error handling

**Handler Layer**:
- Mock services
- Test HTTP responses
- Test status codes

**Example**:
```go
func TestFileService_UploadFile(t *testing.T) {
    // Arrange
    mockFileRepo := &MockFileRepository{}
    mockFolderRepo := &MockFolderRepository{}
    service := usecase.NewFileService(mockFileRepo, mockFolderRepo, ...)
    
    // Act
    result, err := service.UploadFile(...)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Integration Tests (Recommended)

- Test full flow end-to-end
- Use test database/storage
- Test WebSocket communication

## 🔧 Extension Points

### Add Database Support

1. Create `internal/repository/postgres/` package
2. Implement `domain.FileRepository` interface
3. Inject in `main.go` instead of memory repo

```go
// In main.go
db := connectToPostgres()
fileRepo := postgres.NewFileRepository(db)
// Rest stays the same!
```

### Add Cloud Storage

1. Create `internal/storage/cloudinary/` package
2. Add storage interface to domain
3. Inject into FileService

```go
type StorageProvider interface {
    Upload(file io.Reader) (string, error)
    Delete(url string) error
}
```

### Add Email Notifications

1. Create `internal/notification/email/` package
2. Add notification interface
3. Inject into use cases

```go
type Notifier interface {
    NotifyAdmin(message string) error
}
```

## 📊 Benefits of This Architecture

✅ **Testable**: Each layer tested independently  
✅ **Maintainable**: Clear separation of concerns  
✅ **Flexible**: Easy to swap implementations  
✅ **Scalable**: Add features without breaking existing code  
✅ **Clean**: Business logic independent of frameworks  
✅ **Professional**: Industry-standard patterns  

## 🚀 Future Improvements

1. **Persistence**: Add PostgreSQL/MongoDB
2. **Caching**: Add Redis for session management
3. **Queue**: Add message queue for print jobs
4. **Monitoring**: Add logging and metrics
5. **API Docs**: Add Swagger/OpenAPI
6. **Testing**: Add comprehensive test suite
7. **Docker**: Containerize the application
8. **CI/CD**: Automated testing and deployment

---

**This architecture is production-ready and follows Go best practices!**
