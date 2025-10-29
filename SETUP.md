# Quick Setup Guide

## 📋 Prerequisites
- Go 1.23.6+ installed
- Web browser

## ⚡ Quick Start

### 1. Create your `.env` file

Copy the example and customize:
```bash
copy .env.example .env
```

Edit `.env` with your preferred settings:
```env
PORT=8080
HOST=localhost
ADMIN_USERNAME=admin
ADMIN_PASSWORD=YourSecurePassword123!
JWT_SECRET=change-this-to-a-random-secret-key
MAX_FILE_SIZE=10485760
ALLOWED_EXTENSIONS=jpg,jpeg,png,pdf,gif
STORAGE_PATH=./uploads
```

### 2. Install dependencies (Already done!)

```bash
go mod tidy
```

### 3. Run the application

```bash
go run main.go
```

You should see:
```
🚀 Server starting on http://localhost:8080
📁 User upload page: http://localhost:8080
🔐 Admin login: http://localhost:8080/admin
👤 Admin credentials - Username: admin, Password: changeme123
```

## 🎯 Testing the Application

### Test User Upload Flow:

1. Open `http://localhost:8080` in your browser
2. Enter a folder name: "Test Documents"
3. Select a PDF or image file
4. Click "Upload Files"
5. You should see a success message

### Test Admin Dashboard:

1. Open `http://localhost:8080/admin` in another tab
2. Login with your credentials (from .env or default: admin/changeme123)
3. You should see your uploaded files in real-time
4. Try the Print button - it opens the file in a new window
5. Try the Delete button - removes the file instantly

### Test Real-time Updates:

1. Keep admin dashboard open
2. Go back to upload page
3. Upload another file
4. Watch it appear instantly in the admin dashboard! ✨

## 📁 Project Structure

```
fileprintapp/
├── internal/
│   ├── config/          ← Configuration loading
│   ├── domain/          ← Business entities & interfaces
│   ├── handler/         ← HTTP request handlers
│   ├── middleware/      ← Auth, CORS middleware
│   ├── repository/      ← Data storage (in-memory)
│   ├── usecase/         ← Business logic
│   └── websocket/       ← Real-time updates
├── web/static/          ← Frontend files
├── uploads/             ← Uploaded files (auto-created)
├── .env                 ← Your configuration
├── main.go              ← Entry point
└── README.md            ← Full documentation
```

## 🔧 Troubleshooting

### Port already in use?
Change `PORT=8080` to another port in `.env` like `PORT=3000`

### Can't access from other devices?
Change `HOST=localhost` to `HOST=0.0.0.0` in `.env`

### Uploads not working?
- Check file size (default max: 10MB)
- Check file type (only jpg, jpeg, png, pdf, gif allowed)
- Check console for errors

### Admin login not working?
- Verify credentials in `.env` file
- Check browser console for errors
- Clear browser cache/cookies

## 🌟 Features to Try

✅ **Multiple File Upload**: Select multiple files at once  
✅ **Real-time Updates**: Watch files appear instantly  
✅ **Folder Organization**: Group files by folder name  
✅ **Direct Print**: Print without downloading  
✅ **Clean UI**: Modern gradient design  

## 🚀 Next Steps

1. **Change default password** in `.env`
2. **Test on your local network** (use HOST=0.0.0.0)
3. **Deploy to production** (see README.md)
4. **Integrate cloud storage** for production use

## 📝 Important Notes

- **Storage**: Currently uses in-memory storage. Files lost on restart.
- **Security**: Change admin password before deploying!
- **JWT Secret**: Use a strong random string in production
- **File Cleanup**: Manual deletion required (admin dashboard)

## 🎨 Customization Ideas

- Change colors in `web/static/css/style.css`
- Modify max file size in `.env`
- Add more file types in `.env`
- Extend with database storage
- Add user notifications

---

**Need Help?** Check the full README.md for detailed documentation!
