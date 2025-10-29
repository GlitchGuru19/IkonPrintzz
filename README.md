# File Print Service

A modern web application for easy file printing with real-time admin monitoring using WebSockets. Users can upload files without sign-up, and admins can view and print them instantly.

## 🎯 Features

- **No Sign-up Required**: Users can upload files directly without creating accounts
- **Real-time Updates**: Admin dashboard updates instantly via WebSockets
- **Folder Organization**: Users create named folders to organize their uploads
- **Multiple File Types**: Supports PDF, JPG, PNG, and GIF files
- **Admin Authentication**: Secure JWT-based admin login
- **Print-friendly**: Direct print from browser without downloading
- **Clean Architecture**: Well-structured Go code with separation of concerns
- **Modern UI**: Beautiful, responsive interface with gradient design

## 🏗️ Architecture

The project follows **Clean Architecture** principles:

```
fileprintapp/
├── cmd/
│   └── server/          # Alternative server entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Business entities and interfaces
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware (auth, CORS)
│   ├── repository/      # Data storage implementations
│   │   └── memory/      # In-memory repository
│   ├── usecase/         # Business logic / services
│   └── websocket/       # WebSocket hub and client
├── web/
│   └── static/          # Frontend files
│       ├── css/
│       ├── js/
│       └── *.html
├── .env                 # Environment variables (create from .env.example)
├── .env.example         # Example environment configuration
├── go.mod
└── main.go             # Application entry point
```

### Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Defines business entities and repository interfaces
   - No dependencies on other layers

2. **Repository Layer** (`internal/repository/`)
   - Implements data storage
   - Currently uses in-memory storage (easily extendable to database)

3. **Use Case Layer** (`internal/usecase/`)
   - Contains business logic
   - Orchestrates data flow between layers

4. **Handler Layer** (`internal/handler/`)
   - HTTP request handlers
   - Converts HTTP to/from domain objects

5. **Middleware Layer** (`internal/middleware/`)
   - Authentication, CORS, logging, etc.

## 🚀 Getting Started

### Prerequisites

- Go 1.23.6 or higher
- A modern web browser

### Installation

1. **Clone the repository** (or navigate to your project):
   ```bash
   cd d:\GO\GO PROJECTS\fileprintapp
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Create environment file**:
   ```bash
   copy .env.example .env
   ```

4. **Edit `.env` file** with your settings:
   ```env
   PORT=8080
   HOST=localhost
   ADMIN_USERNAME=admin
   ADMIN_PASSWORD=your-secure-password
   JWT_SECRET=your-secret-key-change-this
   ```

### Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

You'll see output like:
```
🚀 Server starting on http://localhost:8080
📁 User upload page: http://localhost:8080
🔐 Admin login: http://localhost:8080/admin
👤 Admin credentials - Username: admin, Password: changeme123
```

## 📱 Usage

### For Users

1. Navigate to `http://localhost:8080`
2. Enter a folder name to organize your files
3. Select one or more files (PDF, images)
4. Click "Upload Files"
5. Your files are ready for printing!

### For Admin

1. Navigate to `http://localhost:8080/admin`
2. Log in with your admin credentials
3. View all uploaded files organized by folder in real-time
4. Click "🖨️ Print" to print directly from browser
5. Click "🗑️ Delete" to remove files after printing

## 🔒 Security

- **JWT Authentication**: Admin routes protected with JWT tokens
- **Password Hashing**: Bcrypt for secure password storage
- **CORS Configuration**: Configurable cross-origin settings
- **File Validation**: Type and size restrictions
- **Environment Variables**: Sensitive data in `.env` file

## 🛠️ Configuration

All configuration is in `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `HOST` | Server host | `localhost` |
| `ADMIN_USERNAME` | Admin username | `admin` |
| `ADMIN_PASSWORD` | Admin password | `changeme123` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key-change-this` |
| `MAX_FILE_SIZE` | Max file size in bytes | `10485760` (10MB) |
| `ALLOWED_EXTENSIONS` | Allowed file types | `jpg,jpeg,png,pdf,gif` |
| `STORAGE_PATH` | Upload directory | `./uploads` |

## 🌐 API Endpoints

### Public Endpoints

- `POST /api/upload` - Upload a file
- `POST /api/folders` - Create a folder
- `POST /api/admin/login` - Admin login

### Protected Endpoints (Require JWT)

- `GET /api/files` - Get all files
- `DELETE /api/files/{id}` - Delete a file
- `GET /api/files/{id}/view` - View/print a file

### WebSocket

- `WS /ws` - Real-time updates for admin dashboard

## 🔄 WebSocket Messages

The admin dashboard receives real-time updates:

```json
{
  "type": "new_file",
  "payload": { "id": "...", "file_name": "...", ... }
}

{
  "type": "file_deleted",
  "payload": { "id": "..." }
}

{
  "type": "folder_created",
  "payload": { "id": "...", "name": "..." }
}
```

## 📦 Dependencies

- `github.com/gorilla/mux` - HTTP router
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading
- `golang.org/x/crypto` - Password hashing

## 🚀 Deployment

### Production Considerations

1. **Change default credentials** in `.env`
2. **Use strong JWT secret**
3. **Configure CORS** for your domain
4. **Set up HTTPS** with reverse proxy (nginx/caddy)
5. **Consider database** instead of in-memory storage
6. **Set up cloud storage** (Cloudinary, S3, etc.)

### Free Hosting Options

- **Fly.io** - Go app hosting
- **Railway** - Easy deployment
- **Render** - Free tier available

For image hosting:
- **Cloudinary** - Free tier with good API
- **imgbb** - Simple free hosting

## 🎨 Frontend Stack

- Pure HTML5
- CSS3 with modern gradients and animations
- Vanilla JavaScript (no frameworks)
- WebSocket API for real-time updates

## 🤝 Contributing

Contributions are welcome! Areas for improvement:

- [ ] Database integration (PostgreSQL, MongoDB)
- [ ] Cloud storage integration
- [ ] User session management
- [ ] Email notifications
- [ ] Print job queue
- [ ] File preview modal
- [ ] Drag-and-drop upload
- [ ] Progress bars
- [ ] Mobile app

## 📝 License

This project is provided as-is for educational and personal use.

## 👨‍💻 Development

### Project Structure Explanation

- **Clean Architecture**: Dependencies point inward (domain has no deps)
- **Repository Pattern**: Easy to swap storage implementations
- **Dependency Injection**: Services receive dependencies via constructors
- **Interface Segregation**: Small, focused interfaces
- **Separation of Concerns**: Each layer has a single responsibility

### Adding New Features

1. Define entities in `internal/domain/`
2. Add repository interface in `internal/domain/repository.go`
3. Implement repository in `internal/repository/`
4. Create use case in `internal/usecase/`
5. Add handler in `internal/handler/`
6. Register route in `main.go`

---

**Happy Printing! 🖨️**
