# Railway Deployment Guide for ERP System Backend

This guide will help you deploy your ERP system microservices to Railway. Since Railway doesn't support Docker Compose, each service needs to be deployed individually.

## Prerequisites

1. **Railway Account**: Sign up at [railway.app](https://railway.app)
2. **Railway CLI** (optional but recommended):
   ```bash
   npm install -g @railway/cli
   railway login
   ```

## Architecture Overview

Your ERP system consists of:
- **6 Microservices**: auth, org, product, inventory, license, subscription
- **3 Infrastructure Services**: MongoDB, Redis, RabbitMQ (via Railway plugins/templates)

## Deployment Strategy

### Phase 1: Create Railway Project

1. **Create a new project** on Railway Dashboard
2. **Link your GitHub repository** to enable automatic deployments

### Phase 2: Deploy Infrastructure Services

#### 2.1 MongoDB
1. Click **"+ New"** → **"Database"** → **"Add MongoDB"**
2. Railway will provision a MongoDB instance
3. Note the connection string from the **"Connect"** tab (format: `mongodb://...`)
4. Save the following variables from MongoDB service:
   - `MONGO_URL` (or similar - Railway auto-generates this)

#### 2.2 Redis
1. Click **"+ New"** → **"Database"** → **"Add Redis"**
2. Railway will provision a Redis instance
3. Note the connection details:
   - `REDIS_URL` or `REDIS_PRIVATE_URL`

#### 2.3 RabbitMQ (Optional - if needed)
Since Railway doesn't have a native RabbitMQ plugin:
- **Option A**: Use [CloudAMQP](https://www.cloudamqp.com/) (free tier available)
- **Option B**: Deploy RabbitMQ as a template from Railway community templates
- **Option C**: Use Railway's custom Docker deployment

For CloudAMQP:
1. Sign up at cloudamqp.com
2. Create a free instance
3. Copy the AMQP URL

### Phase 3: Deploy Microservices

Each service needs to be deployed separately. Railway will detect the `nixpacks.toml` files and build accordingly.

#### 3.1 Deploy Auth Service

1. **Add New Service**: Click **"+ New"** → **"GitHub Repo"** → Select your repo
2. **Configure Root Directory**: 
   - Set **Root Directory** to `/` (since nixpacks.toml handles the service-specific build)
   - Or create separate Railway services pointing to different branches/paths
3. **Add Environment Variables**:
   ```
   PORT=8001
   MONGO_URI=${{MongoDB.MONGO_URL}}
   REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
   JWT_SECRET=<generate-strong-secret>
   JWT_EXPIRY=15m
   REFRESH_TOKEN_EXPIRY=168h
   GOOGLE_CLIENT_ID=<your-google-client-id>
   GOOGLE_CLIENT_SECRET=<your-google-client-secret>
   ENV=production
   ```
   
   **Note**: Use Railway's variable referencing like `${{MongoDB.MONGO_URL}}` to reference other services

4. **Deploy**: Railway will automatically build and deploy

#### 3.2 Deploy Organization Service

1. **Add New Service** from the same repository
2. **Configure Path**: Point to `services/org-service` or use nixpacks.toml
3. **Add Environment Variables**:
   ```
   PORT=8002
   MONGO_URI=${{MongoDB.MONGO_URL}}
   REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
   RABBITMQ_URL=${{RabbitMQ.RABBITMQ_URL}}
   JWT_SECRET=<same-as-auth-service>
   ENV=production
   ```

4. **Deploy**

#### 3.3 Deploy Product Service

1. **Add New Service**
2. **Add Environment Variables**:
   ```
   PORT=8003
   MONGO_URI=${{MongoDB.MONGO_URL}}
   REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
   RABBITMQ_URL=${{RabbitMQ.RABBITMQ_URL}}
   JWT_SECRET=<same-as-auth-service>
   ENV=production
   ```

3. **Deploy**

#### 3.4 Deploy Inventory Service

1. **Add New Service**
2. **Add Environment Variables**:
   ```
   PORT=8006
   MONGO_URI=${{MongoDB.MONGO_URL}}
   REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
   JWT_SECRET=<same-as-auth-service>
   ENV=production
   ```

3. **Deploy**

#### 3.5 Deploy License Service

1. **Add New Service**
2. **Add Environment Variables**:
   ```
   PORT=8004
   MONGO_URI=${{MongoDB.MONGO_URL}}
   REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
   JWT_SECRET=<same-as-auth-service>
   ENV=production
   ```

3. **Deploy**

#### 3.6 Deploy Subscription Service

1. **Add New Service**
2. **Add Environment Variables**:
   ```
   PORT=8005
   MONGO_URI=${{MongoDB.MONGO_URL}}
   REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
   JWT_SECRET=<same-as-auth-service>
   ENV=production
   ```

3. **Deploy**

## Important Railway Considerations

### Multi-Service Deployment from Single Repository

Railway supports multiple services from a single repository. You have two options:

#### Option 1: Single Repo, Multiple Services (Recommended)
1. Create one Railway service per microservice
2. For each service, configure the **Watch Paths** to only trigger builds when specific directories change:
   - Auth Service: `services/auth-service/**`, `shared/**`
   - Org Service: `services/org-service/**`, `shared/**`
   - Product Service: `services/product-service/**`, `shared/**`
   - etc.

3. Set different `nixpacks.toml` configs for each service

#### Option 2: Monorepo with Railway CLI
Use the Railway CLI to deploy specific services:
```bash
# Deploy auth service
cd services/auth-service
railway up

# Deploy org service
cd services/org-service
railway up
```

### Environment Variables Best Practices

1. **Shared Secrets**: Use Railway's **Shared Variables** feature for common values like:
   - `JWT_SECRET`
   - `ENV`
   - `MONGO_URI`
   - `REDIS_ADDR`

2. **Service References**: Use Railway's variable referencing:
   ```
   MONGO_URI=${{MongoDB.MONGO_URL}}
   REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
   ```

3. **Generate Strong Secrets**:
   ```bash
   # Generate JWT secret
   openssl rand -base64 32
   ```

### Networking

Railway automatically provides:
- **Private Networking**: Services within the same project can communicate via private URLs
- **Public Domains**: Each service gets a public domain (e.g., `your-service.up.railway.app`)
- **Custom Domains**: You can add custom domains in the service settings

### Database Initialization

Railway MongoDB instances are empty by default. You'll need to:

1. **Run Migrations**: Create a separate Railway service or use Railway's **"Tasks"** feature:
   ```bash
   railway run bash scripts/migrations/run-migrations.sh
   ```

2. **Seed Data**: Similarly, seed your database:
   ```bash
   railway run bash scripts/seeds/run-seeds.sh
   ```

Or create a one-time deployment service that runs these scripts on startup.

## Service Health Checks

Railway automatically monitors service health. Ensure each service has a health endpoint:

```go
// Add to each service's main.go or router
router.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})
})
```

## Monitoring and Logs

1. **View Logs**: Click on any service → **"Logs"** tab
2. **Metrics**: Railway provides CPU, Memory, and Network metrics
3. **Alerts**: Set up alerts for service failures

## Cost Optimization

Railway's pricing is based on:
- **Compute**: CPU and memory usage
- **Storage**: Database storage
- **Bandwidth**: Network egress

Tips to optimize:
1. Use appropriate instance sizes (don't over-provision)
2. Enable **"Auto-Sleep"** for non-production environments
3. Use Redis for caching to reduce database queries
4. Monitor usage in Railway dashboard

## Deployment Checklist

- [ ] MongoDB database created and accessible
- [ ] Redis instance created
- [ ] RabbitMQ configured (if using)
- [ ] All environment variables set correctly
- [ ] JWT_SECRET is strong and consistent across services
- [ ] Auth Service deployed and healthy
- [ ] Org Service deployed and healthy
- [ ] Product Service deployed and healthy
- [ ] Inventory Service deployed and healthy
- [ ] License Service deployed and healthy
- [ ] Subscription Service deployed and healthy
- [ ] Database migrations run successfully
- [ ] Seed data loaded (if needed)
- [ ] Health checks responding correctly
- [ ] Service-to-service communication working
- [ ] Public endpoints accessible

## Troubleshooting

### Build Failures
1. Check the build logs in Railway dashboard
2. Verify `nixpacks.toml` configuration
3. Ensure Go version compatibility (using Go 1.23)
4. Check that shared dependencies are accessible

### Runtime Errors
1. Verify environment variables are set correctly
2. Check service logs for error messages
3. Ensure MongoDB/Redis connections are correct
4. Verify network connectivity between services

### Database Connection Issues
1. Check MongoDB connection string format
2. Ensure services are in the same Railway project (for private networking)
3. Verify database credentials
4. Check if IP whitelisting is needed (usually not with Railway)

### Service Discovery
If services can't communicate:
1. Use Railway's private networking URLs
2. Check environment variables reference correct services
3. Ensure services are in the same project

## Next Steps

After deployment:
1. Set up monitoring and alerting
2. Configure custom domains
3. Enable SSL certificates (automatic with Railway)
4. Set up CI/CD pipelines
5. Configure backup strategies for databases
6. Implement logging aggregation
7. Set up staging environment

## Additional Resources

- [Railway Documentation](https://docs.railway.app)
- [Nixpacks Documentation](https://nixpacks.com)
- [Railway Discord Community](https://discord.gg/railway)

## Example: Deploying with Railway CLI

```bash
# Initialize Railway project
railway init

# Link to existing project
railway link

# Deploy auth service
cd services/auth-service
railway up

# Set environment variables
railway variables set JWT_SECRET=your-secret-here
railway variables set MONGO_URI=mongodb://...

# View logs
railway logs

# Open service in browser
railway open
```

## Environment Variables Template

Create a `.env.railway` file for each service (don't commit this):

### Auth Service
```env
PORT=8001
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<generate-this>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
ENV=production
```

### Org Service
```env
PORT=8002
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
RABBITMQ_URL=${{RabbitMQ.RABBITMQ_URL}}
JWT_SECRET=<same-as-auth>
ENV=production
```

### Product Service
```env
PORT=8003
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
RABBITMQ_URL=${{RabbitMQ.RABBITMQ_URL}}
JWT_SECRET=<same-as-auth>
ENV=production
```

### Inventory Service
```env
PORT=8006
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<same-as-auth>
ENV=production
```

### License Service
```env
PORT=8004
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<same-as-auth>
ENV=production
```

### Subscription Service
```env
PORT=8005
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<same-as-auth>
ENV=production
```

---

**Note**: This guide assumes you're deploying from a single GitHub repository. Adjust the watch paths and build configurations according to your specific needs.
