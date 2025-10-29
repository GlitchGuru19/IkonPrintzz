# ⚡ QUICK START - Production Deployment

## 🎯 Get Your App Live in 10 Minutes!

---

## Step 1: Create Your .env File (1 min)

```bash
# In your project folder, create .env file with these values:
copy .env.example .env
```

Then edit `.env` with your actual values:

```env
# ====================
# REQUIRED: CHANGE THESE!
# ====================
ADMIN_PASSWORD=Your_Secure_Password_Here_123!
JWT_SECRET=your-random-32-character-secret-key-here-change-this

# ====================
# NEON DATABASE (from your Neon dashboard)
# ====================
DB_HOST=ep-your-project.region.aws.neon.tech
DB_USER=your-neon-username
DB_PASSWORD=your-neon-password
DB_NAME=neondb
DB_PORT=5432
DB_SSL_MODE=require

# ====================
# OPTIONAL: Leave as default
# ====================
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=production
ADMIN_USERNAME=admin
MAX_FILE_SIZE=10485760
ALLOWED_EXTENSIONS=jpg,jpeg,png,pdf,gif
STORAGE_TYPE=local
STORAGE_PATH=./uploads
```

**🔑 Generate JWT Secret:**
```bash
# Windows PowerShell:
-join ((65..90) + (97..122) + (48..57) | Get-Random -Count 32 | % {[char]$_})

# Or just use this online: https://randomkeygen.com/
```

---

## Step 2: Get Your Neon Database Connection (2 min)

1. Go to [Neon Console](https://console.neon.tech/)
2. Click on your project
3. Click **"Connection Details"**
4. Copy the values:
   - **Host**: `ep-xxxxx.region.aws.neon.tech`
   - **Database**: `neondb`
   - **User**: Your username
   - **Password**: Your password
5. Paste these into your `.env` file

**Your Neon database is already set to "production"** ✓

---

## Step 3: Choose Your Hosting Platform (Pick ONE)

### 🚂 OPTION A: Railway (RECOMMENDED - Easiest)

#### Why Railway?
- ✅ Easiest setup
- ✅ $5/month free credit
- ✅ Auto-deploy from GitHub
- ✅ Automatic HTTPS

#### Deploy to Railway:

1. **Push to GitHub** (if not already):
   ```bash
   git add .
   git commit -m "Ready for production"
   git push origin main
   ```

2. **Go to Railway**:
   - Visit [railway.app](https://railway.app/)
   - Click "Start a New Project"
   - Choose "Deploy from GitHub repo"
   - Select your repository

3. **Set Environment Variables**:
   - In Railway dashboard, click "Variables"
   - Click "RAW Editor"
   - Paste ALL variables from your `.env` file
   - Click "Save"

4. **Deploy**:
   - Railway auto-deploys!
   - Wait 1-2 minutes
   - Click the generated URL

5. **Access Your App**:
   - **User Page**: `https://your-app.up.railway.app`
   - **Admin Page**: `https://your-app.up.railway.app/admin` ⚠️ Type this manually!

**DONE!** 🎉

---

### ✈️ OPTION B: Fly.io (For Global Edge Network)

```bash
# Install Fly CLI (PowerShell as Admin):
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"

# Login:
fly auth login

# Launch app:
fly launch

# Set secrets (one by one):
fly secrets set ADMIN_PASSWORD=yourpassword
fly secrets set JWT_SECRET=yoursecret
fly secrets set DB_HOST=your-neon-host
fly secrets set DB_NAME=neondb
fly secrets set DB_USER=your-user
fly secrets set DB_PASSWORD=your-db-password

# Deploy:
fly deploy

# Open:
fly open
```

**Admin Access**: `https://your-app.fly.dev/admin`

---

### 🎨 OPTION C: Render (For Simple Dashboard)

1. Go to [render.com](https://render.com)
2. Click "New +" → "Web Service"
3. Connect GitHub repository
4. Configure:
   - **Name**: fileprintapp
   - **Build Command**: `go build -o bin/server main.go`
   - **Start Command**: `./bin/server`
5. Add Environment Variables (from your .env file)
6. Click "Create Web Service"

**Admin Access**: `https://fileprintapp.onrender.com/admin`

---

## Step 4: Test Your Deployment (2 min)

### Test 1: User Upload

1. Go to your deployed URL
2. Enter folder name: "Test"
3. Upload a file
4. Should see success message ✓

### Test 2: Admin Access

1. Go to `https://your-url.com/admin` (type `/admin` manually!)
2. Login with your credentials
3. Should see uploaded file ✓
4. Green dot = WebSocket connected ✓

### Test 3: Real-Time Update

1. Keep admin dashboard open
2. Open main page in new tab
3. Upload another file
4. Watch it appear instantly in admin! ✓

---

## 🔐 CRITICAL: Admin Access Information

### ⚠️ How to Access Admin Panel When Hosted:

**The admin login link is intentionally HIDDEN from the user page for security!**

**To access admin:**
1. Go to your main URL
2. **Manually type** `/admin` at the end
3. Example: `https://myapp.up.railway.app/admin`

**DO NOT share this URL publicly!**

- ✅ Share with users: `https://myapp.up.railway.app`
- ❌ Keep private: `https://myapp.up.railway.app/admin`

**Bookmark the admin URL** for easy access!

---

## 📊 Monitoring Your App

### Railway Dashboard Shows:

- **Logs**: Real-time application logs
  - Should see: "✅ Database connection established"
  - Should see: "✅ Server is running"
  
- **Deployments**: Each deployment listed
  
- **Metrics**: CPU, RAM, Network usage

### Check if Everything is Working:

```
Look for these in logs:
📋 Loading configuration...
✅ Configuration loaded (Environment: production)
🔌 Connecting to Neon PostgreSQL database...
✅ Database connection established successfully
🔄 Running database migrations...
✅ Database migrations completed successfully
👤 Initializing admin user...
✅ Admin user 'admin' initialized
💾 Initializing repositories...
⚙️  Initializing services...
🔌 Starting WebSocket hub...
🌐 Initializing HTTP handlers...
🛣️  Setting up routes...
✅ Server is running and ready to accept connections!
```

---

## 🎯 Quick Reference

| What | URL |
|------|-----|
| **User Upload Page** | `https://your-app.com` |
| **Admin Login** | `https://your-app.com/admin` |
| **Admin Dashboard** | `https://your-app.com/admin/dashboard` |

**Default Credentials:**
- Username: `admin` (or what you set)
- Password: What you set in .env

---

## ⚠️ Troubleshooting

### "Can't find admin page"
- Type `/admin` manually in URL bar
- Don't look for link on main page (it's hidden!)

### "Database connection failed"
- Check DB_HOST, DB_USER, DB_PASSWORD in Railway variables
- Ensure DB_SSL_MODE=require
- Check Neon dashboard for connection string

### "Login not working"
- Verify ADMIN_PASSWORD in Railway matches what you're typing
- Check JWT_SECRET is set
- Clear browser cookies

### "Files not showing in admin"
- Check green/red dot (WebSocket status)
- Refresh page
- Check browser console (F12)

---

## ✅ Final Checklist

- [ ] .env file created with YOUR values
- [ ] Neon database connection details copied
- [ ] Admin password changed from default
- [ ] JWT secret generated (32+ characters)
- [ ] Pushed to GitHub
- [ ] Deployed to Railway/Fly/Render
- [ ] Environment variables set in hosting platform
- [ ] Tested user upload works
- [ ] Tested admin login at `/admin` URL
- [ ] Tested real-time updates
- [ ] Bookmarked admin URL for future use

---

## 🎉 YOU'RE LIVE!

**Your production-ready file printing service is deployed!**

**Users access**: `https://your-app-url.com`  
**You manage at**: `https://your-app-url.com/admin`

Remember:
- 📝 Files persist in PostgreSQL (won't be lost on restart)
- 🔄 Real-time updates via WebSocket
- 🔐 Admin page hidden from public
- 🚀 Auto-deploy on git push (Railway)

---

## 📱 Next Steps

1. **Custom Domain** (optional):
   - Add in Railway dashboard
   - Point CNAME record to Railway URL
   
2. **Monitor Usage**:
   - Check Railway logs weekly
   - Monitor Neon database size
   - Delete old files after printing

3. **Share with Users**:
   - Give them main URL only
   - Don't mention `/admin`

---

**Questions?** Check PRODUCTION_DEPLOY.md for detailed guides!

**Problems?** Check Railway logs or Neon dashboard!
