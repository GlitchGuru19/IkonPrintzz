# 🚂 Deploy Ikon_Printz to Railway

## ✅ Your app is 100% configured and ready to deploy!

---

## 📋 What's Already Configured

✅ **Neon PostgreSQL Database** - Connected to: `ikondb`  
✅ **JWT Secret** - Set to: `d7d984e723620398426a01a7083952a2`  
✅ **Branding** - All pages show "Ikon_Printz 🖨️"  
✅ **Admin Access** - Hidden from user page  
✅ **Railway Config** - `railway.json` ready  

---

## 🚀 Step-by-Step Deployment (5 Minutes)

### Step 1: Push to GitHub (2 minutes)

```bash
# Initialize git (if not already done)
git init

# Add all files
git add .

# Commit
git commit -m "Ikon_Printz ready for production"

# Create repository on GitHub and push
git remote add origin https://github.com/yourusername/fileprintapp.git
git branch -M main
git push -u origin main
```

### Step 2: Deploy to Railway (2 minutes)

1. **Go to Railway**: [https://railway.app/](https://railway.app/)
2. **Sign up/Login** with GitHub
3. Click **"Start a New Project"**
4. Click **"Deploy from GitHub repo"**
5. Select your `fileprintapp` repository
6. Railway will automatically detect it's a Go app!

### Step 3: Set Environment Variables (1 minute)

1. In Railway dashboard, click on your project
2. Click **"Variables"** tab
3. Click **"RAW Editor"**
4. Copy and paste this ENTIRE block:

```env
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=production
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

5. Click **"Update Variables"**

### Step 4: Wait for Deployment

Railway will:
- Build your Go application
- Run database migrations automatically
- Start the server
- Generate a public URL

**Wait 1-2 minutes for deployment to complete**

---

## 🌐 Access Your Deployed App

### After deployment completes:

1. Click **"Settings"** tab in Railway
2. Scroll down to **"Domains"**
3. You'll see something like: `https://fileprintapp-production.up.railway.app`

**Your URLs:**

- 👥 **User Upload Page**: `https://your-app.up.railway.app`
- 🔐 **Admin Login**: `https://your-app.up.railway.app/admin`

---

## ⚠️ CRITICAL: How to Access Admin Panel

### The admin login is NOT visible on the user page!

**To access admin:**

1. Go to your Railway URL (e.g., `https://fileprintapp-production.up.railway.app`)
2. **Manually type** `/admin` at the end
3. Example: `https://fileprintapp-production.up.railway.app/admin`
4. Login with:
   - **Username**: `admin`
   - **Password**: `changeme123`

**⚠️ BOOKMARK THIS ADMIN URL!**

---

## 📊 Verify Deployment

### Check Railway Logs:

1. Go to Railway dashboard
2. Click **"Deployments"** tab
3. Click on the latest deployment
4. Check logs for these success messages:

```
✅ Configuration loaded (Environment: production)
✅ Database connection established successfully
✅ Database migrations completed successfully
✅ Admin user 'admin' initialized
✅ Server is running and ready to accept connections!
```

### Test User Upload:

1. Visit your Railway URL
2. Enter folder name: "Test"
3. Upload a test file
4. Should see success message ✅

### Test Admin Dashboard:

1. Visit: `https://your-url.up.railway.app/admin`
2. Login with admin/changeme123
3. Should see your uploaded file ✅
4. Green dot = Connected ✅

---

## 🎨 Custom Domain (Optional)

### Add Your Own Domain:

1. In Railway dashboard, go to **"Settings"**
2. Scroll to **"Domains"**
3. Click **"Add Domain"**
4. Enter your domain (e.g., `print.yourdomain.com`)
5. Add CNAME record in your DNS:

```
Type: CNAME
Name: print
Value: your-app.up.railway.app
TTL: 3600
```

**Then access:**
- Users: `https://print.yourdomain.com`
- Admin: `https://print.yourdomain.com/admin`

---

## 🔧 Update Your App

### To deploy updates:

1. Make changes to your code
2. Commit: `git add . && git commit -m "Update"`
3. Push: `git push origin main`
4. Railway **automatically redeploys**! 🎉

---

## 💰 Cost

**Railway Pricing:**
- Free: $5 monthly credit
- Paid: $5/month for Hobby plan

**Neon Database:**
- Free tier: 0.5GB storage
- Your database: `ikondb` on Neon

**Total: ~$5/month** (often covered by free credits!)

---

## 🔐 Security Notes

### Important:

1. ✅ Admin link is hidden from users
2. ✅ Database uses SSL (required)
3. ✅ HTTPS automatic on Railway
4. ⚠️ Change admin password after first login!

### To Change Admin Password:

1. Go to Railway dashboard
2. Click **"Variables"**
3. Change `ADMIN_PASSWORD` value
4. Click **"Update Variables"**
5. App will redeploy automatically

---

## 🎯 Quick Reference

| What | Where |
|------|-------|
| **User Upload** | `https://your-app.up.railway.app` |
| **Admin Login** | `https://your-app.up.railway.app/admin` ⚠️ Type manually! |
| **Default Login** | Username: `admin`, Password: `changeme123` |
| **Railway Dashboard** | [https://railway.app/dashboard](https://railway.app/dashboard) |
| **Neon Dashboard** | [https://console.neon.tech/](https://console.neon.tech/) |

---

## 🆘 Troubleshooting

### "Can't find admin page"
- Type `/admin` manually in URL bar
- Don't look for link on main page

### "Database connection failed"
- Check Railway logs for error
- Verify all DB_* variables are set correctly
- Check Neon dashboard that database is active

### "Login not working"
- Verify ADMIN_PASSWORD in Railway variables
- Try clearing browser cookies
- Check browser console (F12) for errors

### "Files not showing"
- Check green/red connection dot
- Refresh the admin page
- Check Railway logs for errors

---

## ✅ Files You Need to Deploy

These files are already in your repository and ready:

```
✅ main.go - Application entry point
✅ go.mod, go.sum - Dependencies
✅ internal/ - All backend code
✅ web/static/ - Frontend files
✅ railway.json - Railway configuration
✅ .env.example - Configuration template
✅ .gitignore - Excludes sensitive files
```

**You DON'T need:**
- ❌ Docker files (removed)
- ❌ Fly.io config (removed)
- ❌ Render config (removed)
- ❌ .env file (set in Railway dashboard)

---

## 🎉 You're Ready!

### Next Steps:

1. ✅ Push code to GitHub
2. ✅ Deploy on Railway
3. ✅ Set environment variables
4. ✅ Access at `/admin` URL
5. ✅ Test everything works
6. ✅ Change admin password

**Your Ikon_Printz service will be live in 5 minutes!** 🚀

---

## 📱 After Deployment

**Share with users:**
```
Upload your files at: https://your-app.up.railway.app
```

**Keep private (for you only):**
```
Admin dashboard: https://your-app.up.railway.app/admin
```

**Remember**: Admin URL has `/admin` at the end - type it manually!

---

**Questions? Check Railway logs or Neon dashboard!**

**Happy printing! 🖨️**
