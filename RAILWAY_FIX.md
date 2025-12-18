# Railway Service Configuration Guide

## Problem
Nixpacks couldn't detect the build plan because Railway is looking at the root directory instead of individual service directories.

## Solution: Configure Root Directory for Each Service

When creating each service in Railway, you MUST set the **Root Directory** to point to the specific service.

### How to Configure in Railway Dashboard

1. **Create/Select your service** in Railway
2. Go to **Settings** tab
3. Find **"Root Directory"** or **"Source"** section
4. Set the root directory to the service path

### Root Directory Settings for Each Service

#### Auth Service
```
Root Directory: services/auth-service
Watch Paths: services/auth-service/**, shared/**
```

#### Organization Service
```
Root Directory: services/org-service
Watch Paths: services/org-service/**, shared/**
```

#### Product Service
```
Root Directory: services/product-service
Watch Paths: services/product-service/**, shared/**
```

#### Inventory Service
```
Root Directory: services/inventory-service
Watch Paths: services/inventory-service/**, shared/**
```

#### License Service
```
Root Directory: services/license-service
Watch Paths: services/license-service/**, shared/**
```

#### Subscription Service
```
Root Directory: services/subscription-service
Watch Paths: services/subscription-service/**, shared/**
```

## Alternative: Using Railway CLI

Deploy each service individually using the CLI:

```bash
# Deploy Auth Service
cd services/auth-service
railway up --service auth-service

# Deploy Org Service
cd ../org-service
railway up --service org-service

# And so on...
```

## Why This Happens

Railway/Nixpacks scans the root directory for:
- `go.mod` file (for Go projects)
- `package.json` (for Node projects)
- `requirements.txt` (for Python projects)
- etc.

Since your root doesn't have these files (they're in service subdirectories), Nixpacks can't auto-detect the language.

## Quick Fix Steps

1. Go to your Railway dashboard
2. For the failing service, click on it
3. Go to **Settings** → **Source**
4. Set **Root Directory** to the appropriate path (e.g., `services/auth-service`)
5. Click **Deploy** → **Redeploy**

The build should now succeed!
