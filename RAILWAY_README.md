# Railway Deployment - Quick Reference

This directory contains all files needed for Railway deployment.

## 📁 Files Created

### Configuration Files
- **`railway.toml`** - Main Railway configuration
- **`railway.json`** - Alternative JSON config format
- **`.railwayignore`** - Files to exclude from deployment

### Service-Specific Configs (in each service directory)
- **`services/*/nixpacks.toml`** - Build configuration for each service
  - `auth-service/nixpacks.toml`
  - `org-service/nixpacks.toml`
  - `product-service/nixpacks.toml`
  - `inventory-service/nixpacks.toml`
  - `license-service/nixpacks.toml`
  - `subscription-service/nixpacks.toml`

### Documentation
- **`RAILWAY_DEPLOYMENT.md`** - Complete deployment guide
- **`ENVIRONMENT_VARIABLES.md`** - All environment variables documented

### Scripts
- **`scripts/railway-setup.sh`** - Helper script to generate secrets

## 🚀 Quick Start

### 1. Prerequisites
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login to Railway
railway login
```

### 2. Generate Secrets
```bash
# Run the setup script
bash scripts/railway-setup.sh
```

### 3. Create Railway Project
1. Go to [railway.app](https://railway.app)
2. Create a new project
3. Link your GitHub repository

### 4. Add Databases
In Railway dashboard:
- Add MongoDB: **+ New** → **Database** → **MongoDB**
- Add Redis: **+ New** → **Database** → **Redis**

### 5. Deploy Services

For each service, you'll need to:
1. Add a new service from your GitHub repo
2. Set the watch paths (to trigger builds only when that service changes)
3. Add environment variables (see `ENVIRONMENT_VARIABLES.md`)

Example watch paths for auth-service:
```
services/auth-service/**
shared/**
```

## 🔧 Configuration Overview

### Build Process (Nixpacks)
Each service's `nixpacks.toml` configures:
1. **Setup Phase**: Install Go 1.23
2. **Install Phase**: Download dependencies, cache shared module
3. **Build Phase**: Compile Go binary
4. **Start Phase**: Run the service

### Environment Variables
Railway uses variable referencing for service discovery:
```env
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
```

## 📊 Service Ports

| Service      | Port | Health Check |
|--------------|------|--------------|
| Auth         | 8001 | `/health`    |
| Organization | 8002 | `/health`    |
| Product      | 8003 | `/health`    |
| License      | 8004 | `/health`    |
| Subscription | 8005 | `/health`    |
| Inventory    | 8006 | `/health`    |

## 🔐 Required Environment Variables

All services need:
- `PORT` - Service port (unique per service)
- `MONGO_URI` - MongoDB connection string
- `REDIS_ADDR` - Redis connection address
- `JWT_SECRET` - Shared secret for JWT tokens
- `ENV` - Environment (production/staging/development)

See `ENVIRONMENT_VARIABLES.md` for complete list.

## 🛠️ Deployment Methods

### Method 1: Railway Dashboard (Recommended)
1. Create project
2. Add services via UI
3. Configure each service
4. Deploy automatically on push

### Method 2: Railway CLI
```bash
# Deploy specific service
cd services/auth-service
railway up

# Set environment variables
railway variables set JWT_SECRET=xxx
```

### Method 3: GitHub Actions (Advanced)
Set up CI/CD pipeline with Railway's GitHub integration.

## 📝 Deployment Checklist

- [ ] Railway account created
- [ ] Project initialized
- [ ] GitHub repo linked
- [ ] MongoDB added to project
- [ ] Redis added to project
- [ ] RabbitMQ configured (CloudAMQP or alternative)
- [ ] Secrets generated (JWT_SECRET)
- [ ] Auth service deployed
- [ ] Organization service deployed
- [ ] Product service deployed
- [ ] Inventory service deployed
- [ ] License service deployed
- [ ] Subscription service deployed
- [ ] Environment variables set for all services
- [ ] Health checks verified
- [ ] Database migrations run
- [ ] Seed data loaded (optional)

## 🐛 Troubleshooting

### Build Fails
- Check `nixpacks.toml` configuration
- Verify Go module dependencies
- Check Railway build logs

### Service Won't Start
- Verify environment variables
- Check service logs in Railway dashboard
- Ensure MongoDB/Redis are running

### Services Can't Communicate
- Verify all services in same Railway project
- Use Railway's private URLs
- Check environment variable references

## 📚 Additional Resources

- [RAILWAY_DEPLOYMENT.md](../RAILWAY_DEPLOYMENT.md) - Detailed guide
- [ENVIRONMENT_VARIABLES.md](../ENVIRONMENT_VARIABLES.md) - All variables
- [Railway Docs](https://docs.railway.app)
- [Nixpacks Docs](https://nixpacks.com)

## 💡 Pro Tips

1. **Use Shared Variables**: Set common variables (like JWT_SECRET) once
2. **Watch Paths**: Configure to only rebuild when relevant files change
3. **Private Networking**: Use Railway's internal URLs for service communication
4. **Staging Environment**: Create separate Railway project for staging
5. **Monitoring**: Use Railway's built-in metrics and logs

## 🆘 Getting Help

- Railway Discord: [discord.gg/railway](https://discord.gg/railway)
- Railway Docs: [docs.railway.app](https://docs.railway.app)
- GitHub Issues: Create an issue in your repo

---

**Ready to deploy?** Start with `RAILWAY_DEPLOYMENT.md` for the complete guide!
