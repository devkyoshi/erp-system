# 🚂 Railway Deployment - Complete Solution

## ✅ What's Been Set Up

Your repository is now **Railway-ready** with a clean monorepo deployment solution!

### 📦 Files Created/Modified

1. **Deployment Scripts** (Interactive & Automated)
   - `scripts/railway-init.sh` - One-time project setup
   - `scripts/deploy-railway.sh` - Smart deployment with change detection

2. **Configuration Files**
   - `railway-services.json` - Service definitions for monorepo
   - `railway.toml` - Railway build settings
   - `services/*/nixpacks.toml` - Build config for each service (kept simple)

3. **Documentation**
   - `RAILWAY_SIMPLE_GUIDE.md` - **START HERE** (quick 3-step guide)
   - `RAILWAY_DEPLOYMENT.md` - Detailed documentation
   - `ENVIRONMENT_VARIABLES.md` - All env vars explained

4. **Preserved**
   - All `Dockerfile`s - Kept for Docker Compose and future deployments
   - Existing project structure - No changes to your code

### 🎯 The Solution

**Monorepo Approach**: Deploy all services from a single repository with smart detection.

```
┌─────────────────────────────────────────┐
│   Your GitHub Repository (Main Source)  │
│   • Dockerfiles (for local dev)        │
│   • nixpacks.toml (for Railway)        │
│   • Deployment scripts                  │
└──────────────────┬──────────────────────┘
                   │
                   │ Railway watches specific paths
                   ▼
┌─────────────────────────────────────────┐
│         Railway Project                 │
├─────────────────────────────────────────┤
│  Service: auth-service                  │
│  Watches: services/auth-service/**      │
│           shared/**                     │
├─────────────────────────────────────────┤
│  Service: org-service                   │
│  Watches: services/org-service/**       │
│           shared/**                     │
├─────────────────────────────────────────┤
│  ... (4 more services)                  │
└─────────────────────────────────────────┘
```

## 🚀 How to Deploy (3 Steps)

### 1️⃣ Initial Setup (Run Once)
```bash
bash scripts/railway-init.sh
```
This creates your Railway project and generates secrets.

### 2️⃣ Configure Railway Dashboard
- Create 6 services (one per microservice)
- Add MongoDB and Redis databases
- Set environment variables
- Configure watch paths

**Detailed instructions**: See `RAILWAY_SIMPLE_GUIDE.md`

### 3️⃣ Deploy
```bash
bash scripts/deploy-railway.sh
```
Choose:
- Option 1: Deploy all services
- Option 2: Deploy specific service
- Option 3: Auto-detect and deploy only changed services

## 🎨 Key Features

### ✨ Smart Deployment
The deployment script automatically detects which services changed:
```bash
# Only auth-service and shared changed?
# → Deploy only auth-service

# Shared module changed?
# → Deploy ALL services (they all depend on it)
```

### 🔄 Dual Build System
- **Local Development**: Use Docker Compose with Dockerfiles
- **Railway Production**: Use Nixpacks with nixpacks.toml
- Both coexist peacefully!

### 📊 Watch Paths
Railway only rebuilds services when their specific files change:
```
services/auth-service/** → auth-service rebuilds
services/org-service/** → org-service rebuilds  
shared/** → ALL services rebuild
```

### 🔐 Secure Secrets
- Generated automatically during setup
- Stored in `.env.railway` (gitignored)
- Never committed to repository

## 📋 Deployment Workflow

```bash
# Step 1: Make changes to your code
vim services/auth-service/internal/handlers/auth_handler.go

# Step 2: Commit and push
git add .
git commit -m "Update auth handler"
git push

# Step 3: Deploy (script detects changes)
bash scripts/deploy-railway.sh
# Choose option 3 (smart deploy)
# → Only auth-service deploys!
```

## 🆚 Comparison: Docker vs Railway

| Aspect | Docker Compose (Local) | Railway (Production) |
|--------|----------------------|---------------------|
| **Config File** | `docker-compose.yml` | `nixpacks.toml` |
| **Build** | `Dockerfile` | Nixpacks auto-build |
| **Network** | Docker network | Railway private network |
| **Deploy** | `docker-compose up` | `bash scripts/deploy-railway.sh` |
| **Env Vars** | `.env` file | Railway dashboard |
| **Databases** | Local containers | Railway managed DBs |

## 🔧 Quick Commands Reference

```bash
# Setup (once)
bash scripts/railway-init.sh

# Deploy all
bash scripts/deploy-railway.sh  # → option 1

# Deploy specific service
bash scripts/deploy-railway.sh  # → option 2

# Smart deploy (changed only)
bash scripts/deploy-railway.sh  # → option 3

# Check status
railway status

# View logs
railway logs --service auth-service

# Open dashboard
railway open
```

## 📚 Documentation Files (Read in Order)

1. **`RAILWAY_SIMPLE_GUIDE.md`** ← **START HERE**
   Quick 3-step guide to get deployed fast

2. **`ENVIRONMENT_VARIABLES.md`**
   Complete list of all environment variables

3. **`RAILWAY_DEPLOYMENT.md`**
   Detailed deployment guide with troubleshooting

## ⚙️ For Each Service

Railway configuration per service (repeat 6 times):

| Setting | Value |
|---------|-------|
| **Source** | GitHub repo: `erp-system` |
| **Root Directory** | `/` (repository root) |
| **Watch Paths** | `services/<service-name>/**`<br>`shared/**` |
| **Build** | Auto (uses nixpacks.toml) |
| **Start Command** | Auto (from nixpacks.toml) |

## 🎯 Why This Solution is Clean

1. **No Code Changes**: Your services work exactly as before
2. **Keeps Dockerfiles**: Still use Docker Compose locally
3. **Automated Deployment**: One command deploys everything
4. **Smart Detection**: Only changed services deploy
5. **Monorepo Friendly**: All services in one repo
6. **Production Ready**: Proper environment separation

## ✅ Pre-Deployment Checklist

Before you start:
- [ ] Railway account created
- [ ] Railway CLI installed (`npm install -g @railway/cli`)
- [ ] Git repository pushed to GitHub
- [ ] Node.js installed (for Railway CLI)

## 🐛 Common Issues & Solutions

### "Railway CLI not found"
```bash
npm install -g @railway/cli
```

### "Can't find shared module"
- Make sure Root Directory is `/` (not the service directory)

### "All services rebuild every time"
- Check watch paths are configured correctly per service

### "Build succeeds but service won't start"
- Check environment variables in Railway dashboard
- Verify MongoDB and Redis are added to project

## 💰 Estimated Costs

Railway pricing (as of 2025):
- **Free tier**: $5 credit/month
- **Compute**: ~$0.000463/GB-hour
- **Estimated monthly**: $40-70 for all services

**Cost optimization**:
- Use watch paths to avoid unnecessary rebuilds
- Monitor resource usage in Railway dashboard
- Use staging environment with auto-sleep

## 🎉 You're All Set!

Everything is configured and ready to deploy. Just run:

```bash
# 1. Setup
bash scripts/railway-init.sh

# 2. Configure in Railway dashboard (see RAILWAY_SIMPLE_GUIDE.md)

# 3. Deploy
bash scripts/deploy-railway.sh
```

**Questions?** Check `RAILWAY_SIMPLE_GUIDE.md` for step-by-step instructions!

---

**Happy deploying! 🚀**
