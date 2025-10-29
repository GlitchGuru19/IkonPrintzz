# 🚀 Production Deployment Guide

## Complete Guide for Deploying File Print Service to Production

---

## 📋 Prerequisites Checklist

- ✅ Neon PostgreSQL database created and configured
- ✅ Admin credentials chosen (strong password!)
- ✅ JWT secret generated (32+ characters)
- ✅ Hosting platform account (Railway recommended)
- ✅ Domain name (optional but recommended)

---

## 🎯 RECOMMENDED HOSTING: Railway

**Why Railway?**
- ✅ **FREE tier** available
- ✅ Built-in PostgreSQL support (or use external Neon)
- ✅ Automatic HTTPS
- ✅ Easy environment variable management
- ✅ Git-based deployments
- ✅ Custom domains supported

### Railway Deployment Steps

#### 1. Prepare Your Neon Database

1. Go to [Neon Console](https://console.neon.tech/)
2. Create a new project (if you haven't already)
3. Copy your connection details:
   ```
   Host: ep-xxxxx.region.aws.neon.tech
   Database: neondb
   User: your-username
   Password: your-password
   Port: 5432
   ```

#### 2. Create Railway Project

1. Go to [Railway.app](https://railway.app/)
2. Click **"Start a New Project"**
3. Choose **"Deploy from GitHub repo"**
4. Connect your GitHub account
5. Select your repository

#### 3. Configure Environment Variables

In Railway dashboard, go to **"Variables"** and add ALL of these:

```env
# Server
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=production

# Admin (⚠️ CHANGE THESE!)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=YourSecurePassword123!

# JWT Secret (⚠️ GENERATE NEW ONE!)
JWT_SECRET=your-random-32-character-secret-key-here

# Database (from Neon)
DB_HOST=ep-xxxxx.region.aws.neon.tech
DB_PORT=5432
DB_NAME=neondb
DB_USER=your-neon-username
DB_PASSWORD=your-neon-password
DB_SSL_MODE=require

# Files
MAX_FILE_SIZE=10485760
ALLOWED_EXTENSIONS=jpg,jpeg,png,pdf,gif
STORAGE_TYPE=local
STORAGE_PATH=./uploads
```

#### 4. Deploy

1. Railway will automatically build and deploy
2. Wait for deployment (usually 1-2 minutes)
3. Click on the generated URL (e.g., `https://fileprintapp-production.up.railway.app`)

#### 5. Access Admin Panel

Your app is now live!

- **User Page**: `https://your-railway-app.up.railway.app`
- **Admin Login**: `https://your-railway-app.up.railway.app/admin`

⚠️ **IMPORTANT**: The admin login link is NOT shown on the user page for security.
You must manually type `/admin` at the end of your URL.

---

## 🌐 Alternative Hosting Options

### Option 2: Fly.io

**Pros**: Free tier, global edge network, custom domains
**Setup**: 
1. Install Fly CLI: `powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"`
2. Login: `fly auth login`
3. Create app: `fly launch`
4. Set secrets: `fly secrets set ADMIN_PASSWORD=yourpass`
5. Deploy: `fly deploy`

[Full Fly.io Guide](https://fly.io/docs/languages-and-frameworks/golang/)

### Option 3: Render

**Pros**: Simple setup, free tier, auto-deploy from Git
**Setup**:
1. Go to [Render.com](https://render.com)
2. Create new **Web Service**
3. Connect GitHub repo
4. Set environment variables in dashboard
5. Deploy

[Full Render Guide](https://render.com/docs)

### Option 4: Heroku

**Pros**: Mature platform, add-ons marketplace
**Note**: No longer has free tier, but affordable ($7/month)
**Setup**:
1. Install Heroku CLI
2. `heroku create your-app-name`
3. `heroku config:set` for each env var
4. `git push heroku main`

---

## 🔐 Security Checklist for Production

### Before Deploying:

- [ ] Changed default admin password
- [ ] Generated strong JWT secret (use `openssl rand -base64 32`)
- [ ] Set ENVIRONMENT=production
- [ ] Verified DB_SSL_MODE=require for Neon
- [ ] Tested admin login works
- [ ] Tested file upload works
- [ ] Verified WebSocket connections work

### After Deploying:

- [ ] Admin login link NOT visible on user page ✓
- [ ] HTTPS enabled (automatic on Railway/Fly/Render)
- [ ] Admin password changed from default
- [ ] Database connection secure (SSL enabled)
- [ ] Test file upload and deletion
- [ ] Test real-time updates in admin dashboard

---

## 🗄️ Database Setup (Neon PostgreSQL)

### Your Neon database should already be set to "production" environment.

**What happens on first deployment:**

1. App connects to Neon database
2. Runs migrations automatically (creates tables)
3. Creates admin user with your credentials
4. Ready to use!

**Tables created:**
- `folders` - User-created folders
- `uploaded_files` - File metadata
- `admins` - Admin credentials

**No manual SQL needed** - migrations run automatically!

---

## 🌐 Custom Domain Setup (Optional)

### Railway Custom Domain

1. In Railway dashboard, go to **"Settings"**
2. Click **"Add Domain"**
3. Enter your domain (e.g., `print.yourdomain.com`)
4. Add CNAME record in your DNS:
   ```
   Type: CNAME
   Name: print
   Value: your-app.up.railway.app
   ```
5. Wait for DNS propagation (5-30 minutes)

**Admin Access with Custom Domain:**
- User Page: `https://print.yourdomain.com`
- Admin Page: `https://print.yourdomain.com/admin`

---

## 📊 Monitoring Your App

### Railway Dashboard Shows:

- **Logs**: Real-time application logs
- **Metrics**: CPU, Memory, Network usage
- **Deployments**: Deployment history
- **Variables**: Environment variables (secure)

### Watch for:

- Database connection errors
- File upload failures
- Memory usage (check if uploads directory growing)
- WebSocket connection issues

---

## 🔄 Updating Your App

### Railway Auto-Deploy:

1. Make changes to code
2. Push to GitHub: `git push origin main`
3. Railway automatically detects and redeploys
4. Zero downtime deployment!

### Manual Redeploy:

In Railway dashboard, click **"Redeploy"**

---

## 📱 Testing Your Production App

### 1. Test User Upload

1. Go to your deployed URL
2. Create a folder: "Test Folder"
3. Upload a PDF or image
4. Should see success message

### 2. Test Admin Dashboard

1. Go to `https://your-url.com/admin`
2. Login with your credentials
3. Should see your uploaded file in real-time
4. Click "Print" - should open file
5. Click "Delete" - should remove file

### 3. Test WebSocket

1. Keep admin dashboard open
2. In another tab, upload a new file
3. Watch it appear instantly in admin dashboard
4. Green dot should show "connected"

---

## ⚠️ Troubleshooting Common Issues

### "Database connection failed"

**Fix**: Verify these in Railway environment variables:
- DB_HOST is correct from Neon
- DB_PASSWORD is correct
- DB_SSL_MODE is "require"

### "Admin login not working"

**Fix**: 
- Ensure ADMIN_PASSWORD matches what you're entering
- Check JWT_SECRET is set
- Clear browser cookies
- Check browser console for errors

### "Files not appearing in admin dashboard"

**Fix**:
- Check WebSocket connection (green dot in dashboard)
- Refresh the page
- Check browser console for WebSocket errors
- Verify API calls in Network tab

### "Uploads failing"

**Fix**:
- Check MAX_FILE_SIZE (default 10MB)
- Verify file extension is allowed
- Check STORAGE_PATH directory permissions
- Review application logs in Railway

---

## 💰 Cost Estimate

### Recommended Setup:

| Service | Plan | Cost |
|---------|------|------|
| Railway | Starter | $5/month (includes $5 credit) |
| Neon PostgreSQL | Free | $0 (0.5GB storage) |
| **Total** | | **~$5/month** |

**Railway free trial**: $5 credit monthly (often enough for small apps!)

---

## 🎓 Post-Deployment Best Practices

1. **Monitor Database Size**: Neon free tier has 0.5GB limit
2. **Clean Up Old Files**: Manually delete files after printing
3. **Backup Database**: Neon has automatic backups
4. **Update Dependencies**: Run `go get -u` monthly
5. **Check Logs**: Review Railway logs weekly
6. **Test Admin Access**: Login weekly to verify functionality

---

## 📞 Support Resources

- **Railway Docs**: https://docs.railway.app/
- **Neon Docs**: https://neon.tech/docs/
- **Go Deployment**: https://go.dev/doc/install
- **This App's GitHub**: [Your repo URL]

---

## ✅ Final Checklist

Before going live:

- [ ] Environment variables all set in Railway
- [ ] Neon database connected and migrated
- [ ] Admin password changed from default
- [ ] JWT secret is random and secure
- [ ] Test upload works
- [ ] Test admin login works
- [ ] Test real-time WebSocket updates
- [ ] HTTPS working (automatic on Railway)
- [ ] `/admin` URL manually tested and working
- [ ] User page does NOT show admin link ✓

---

## 🎉 You're Live!

Your production-ready printing service is now deployed!

**Share with users**: `https://your-app.up.railway.app`  
**Keep private**: `https://your-app.up.railway.app/admin`

**Remember**: Users upload files at main URL, you manage them at `/admin`!

---

**Questions? Check logs in Railway dashboard or review code comments in main.go**
