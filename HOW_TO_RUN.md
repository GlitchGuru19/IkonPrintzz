# 🚀 How to Run Ikon_Printz

## Two Options: Local Testing or Production Deploy

---

## 🖥️ Option 1: Run Locally (Testing)

### Quick Local Test (Without Database)

If you just want to test the app works before deploying:

```bash
# 1. Make sure you're in the project directory
cd "d:\GO\GO PROJECTS\fileprintapp"

# 2. Install dependencies
go mod tidy

# 3. Run the app
go run main.go
```

**⚠️ This will fail with database error** because you haven't set up environment variables locally.

### Proper Local Setup (With Database)

If you want to run locally with the real database:

**Step 1: Create local .env file**

Create a file named `.env` in the project root:

```env
PORT=8080
HOST=localhost
ENVIRONMENT=development
ADMIN_USERNAME=admin
ADMIN_PASSWORD=changeme123
JWT_SECRET=d7d984e723620398426a01a7083952a2
DB_HOST=ep-dry-resonance-ah9jtfim-pooler.c-3.us-east-1.aws.neon.tech
DB_PORT=5432
DB_NAME=ikondb
DB_USER=neondb_owner
DB_PASSWORD=npg_3QzZibaWpG1M
DB_SSL_MODE=require
MAX_FILE_SIZE=10485760
ALLOWED_EXTENSIONS=jpg,jpeg,png,pdf,gif
STORAGE_TYPE=local
STORAGE_PATH=./uploads
```

**Step 2: Run the app**

```bash
go run main.go
```

**Step 3: Access locally**

- **User page**: http://localhost:8080
- **Admin page**: http://localhost:8080/admin (type manually!)

**Login:** admin / changeme123

---

## 🌐 Option 2: Deploy to Railway (RECOMMENDED)

### This is the proper way to use your app in production!

**Follow the complete guide:** `DEPLOY_TO_RAILWAY.md`

**Quick Summary:**

1. **Push to GitHub**
   ```bash
   git add .
   git commit -m "Ready for deployment"
   git push origin main
   ```

2. **Go to Railway**: [https://railway.app/](https://railway.app/)

3. **Deploy from GitHub repo**

4. **Set environment variables** (copy from .env.example)

5. **Access your app:**
   - User: `https://your-app.up.railway.app`
   - Admin: `https://your-app.up.railway.app/admin`

**Deployment time: 5 minutes**

---

## ⚠️ Important: Admin Access

### The admin login is HIDDEN from the user page!

**To access admin, you must type `/admin` manually in the URL:**

- ❌ No link on the user page
- ✅ Type `/admin` in browser address bar
- ✅ Bookmark the admin URL

---

## 🔍 Verify Everything Works

### Check Startup Logs

When you run the app, you should see:

```
📋 Loading configuration...
✅ Configuration loaded (Environment: production)
🔌 Connecting to Neon PostgreSQL database...
✅ Database connection established successfully
🔄 Running database migrations...
✅ Database migrations completed successfully
👤 Initializing admin user...
✅ Admin user 'admin' initialized
📁 Setting up file storage...
💾 Initializing repositories...
⚙️  Initializing services...
🔌 Starting WebSocket hub...
🌐 Initializing HTTP handlers...
🛣️  Setting up routes...

============================================================
🚀 Ikon_Printz - PRODUCTION MODE
============================================================
📍 Server Address: http://0.0.0.0:8080
📁 User Upload Page: http://0.0.0.0:8080
🔐 Admin Login: http://0.0.0.0:8080/admin

⚠️  ADMIN ACCESS INSTRUCTIONS:
   When hosted, access admin at: https://yourdomain.com/admin
   The admin link is NOT shown on the user page for security
   Username: admin
============================================================

✅ Server is running and ready to accept connections!
```

**All steps should show ✅ checkmarks!**

---

## 🧪 Test Your App

### Test 1: User Upload

1. Go to your URL (local or Railway)
2. Enter folder name: "Test Folder"
3. Select a PDF or image
4. Click "Upload Files"
5. Should see: ✅ "Files uploaded successfully!"

### Test 2: Admin Login

1. Go to URL + `/admin` (type manually!)
2. Login: admin / changeme123
3. Should redirect to dashboard
4. Should see "Ikon_Printz Dashboard" header

### Test 3: Real-Time Updates

1. Keep admin dashboard open
2. Open user page in another tab
3. Upload a file
4. Watch it appear instantly in admin! ✨
5. Green dot = Connected

### Test 4: Print & Delete

1. In admin dashboard, click "🖨️ Print"
2. File should open in new tab
3. Use browser print (Ctrl+P)
4. Click "🗑️ Delete"
5. File should disappear

**All tests pass? You're good to go!** ✅

---

## 📁 Project Structure

```
fileprintapp/
├── main.go              ← Start here
├── go.mod, go.sum       ← Dependencies
├── .env.example         ← Configuration template
├── .env                 ← Your local config (create this)
├── railway.json         ← Railway deployment
│
├── internal/            ← Backend code
│   ├── config/          ← Load .env variables
│   ├── database/        ← PostgreSQL connection
│   ├── domain/          ← Business models
│   ├── handler/         ← HTTP endpoints
│   ├── middleware/      ← Auth & CORS
│   ├── repository/      ← Database operations
│   ├── usecase/         ← Business logic
│   └── websocket/       ← Real-time updates
│
├── web/static/          ← Frontend
│   ├── index.html       ← User upload page
│   ├── admin-login.html ← Admin login
│   ├── admin-dashboard.html ← Admin dashboard
│   ├── css/style.css    ← Styling
│   └── js/              ← JavaScript
│
└── uploads/             ← Uploaded files (created automatically)
```

---

## 🔧 Development Commands

### Install Dependencies
```bash
go mod tidy
```

### Run Locally
```bash
go run main.go
```

### Build Binary
```bash
go build -o ikonprintz.exe main.go
```

### Run Binary
```bash
./ikonprintz.exe
```

### Check for Errors
```bash
go vet ./...
```

---

## 💡 Quick Tips

### For Local Development:
- Use `HOST=localhost` in .env
- Use `ENVIRONMENT=development`
- Keep Railway for production

### For Production (Railway):
- Use `HOST=0.0.0.0` (allows external access)
- Use `ENVIRONMENT=production`
- Set all variables in Railway dashboard
- Don't commit .env file to Git!

### Database:
- Your Neon database: `ikondb`
- Tables created automatically on first run
- Admin user created automatically
- No manual SQL needed!

---

## 🆘 Common Issues

### "Cannot connect to database"
**Fix:** Check your DB_HOST, DB_USER, DB_PASSWORD in .env

### "Port already in use"
**Fix:** Change PORT=8080 to PORT=3000 (or any free port)

### "Admin page not found"
**Fix:** Type `/admin` manually in URL - it's hidden on purpose!

### "Login not working"
**Fix:** 
- Check ADMIN_PASSWORD in .env
- Try: admin / changeme123
- Clear browser cookies

### "Files not uploading"
**Fix:**
- Check file size (max 10MB)
- Check file type (only jpg, png, pdf, gif)
- Check browser console for errors

---

## ✅ Recommended: Use Railway

**Local testing is fine, but for actual use:**

1. ✅ Deploy to Railway (5 minutes)
2. ✅ Get HTTPS automatically
3. ✅ Access from any device
4. ✅ Reliable 24/7 uptime
5. ✅ Professional hosting

**See: `DEPLOY_TO_RAILWAY.md` for step-by-step guide**

---

## 🎯 Summary

**For Quick Test:**
```bash
go run main.go
```
Access: http://localhost:8080

**For Production:**
1. Deploy to Railway
2. Set environment variables
3. Access at your Railway URL
4. Admin at: `/admin` (type manually!)

---

## 📞 Need Help?

- **Deployment:** Read `DEPLOY_TO_RAILWAY.md`
- **Database:** Check Neon dashboard
- **Logs:** Check Railway dashboard → Logs tab
- **Admin:** Type `/admin` in URL manually!

---

**Your Ikon_Printz app is ready to run! 🖨️**

**Best way: Deploy to Railway in 5 minutes!**
