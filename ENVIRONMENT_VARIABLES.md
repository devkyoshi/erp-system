# Environment Variables for Railway Deployment

This document lists all required environment variables for each microservice in the ERP system when deploying to Railway.

## Shared Variables (All Services)

These variables should be set consistently across all services:

```env
JWT_SECRET=<generate-strong-random-string>  # Use: openssl rand -base64 32
ENV=production
```

## Infrastructure Services

### MongoDB
Railway automatically provides:
- `MONGO_URL` or `MONGODB_URL` - Full connection string

### Redis
Railway automatically provides:
- `REDIS_URL` or `REDIS_PRIVATE_URL` - Connection URL
- Use the private URL for internal service communication

### RabbitMQ (if using CloudAMQP or similar)
- `RABBITMQ_URL` - AMQP connection string (e.g., `amqp://user:pass@host:5672/`)

## Service-Specific Variables

### 1. Auth Service (`services/auth-service`)

**Port**: 8001

```env
# Required
PORT=8001
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<shared-secret>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
ENV=production

# Optional - Google OAuth
GOOGLE_CLIENT_ID=<your-google-client-id>
GOOGLE_CLIENT_SECRET=<your-google-client-secret>
GOOGLE_REDIRECT_URL=https://your-auth-service.railway.app/auth/google/callback
```

**Health Check**: `/health`

---

### 2. Organization Service (`services/org-service`)

**Port**: 8002

```env
# Required
PORT=8002
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<shared-secret>
ENV=production

# Optional - if using message queue
RABBITMQ_URL=${{RabbitMQ.RABBITMQ_URL}}
```

**Health Check**: `/health`

---

### 3. Product Service (`services/product-service`)

**Port**: 8003

```env
# Required
PORT=8003
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<shared-secret>
ENV=production

# Optional - if using message queue
RABBITMQ_URL=${{RabbitMQ.RABBITMQ_URL}}
```

**Health Check**: `/health`

---

### 4. License Service (`services/license-service`)

**Port**: 8004

```env
# Required
PORT=8004
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<shared-secret>
ENV=production
```

**Health Check**: `/health`

---

### 5. Subscription Service (`services/subscription-service`)

**Port**: 8005

```env
# Required
PORT=8005
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<shared-secret>
ENV=production

# Optional - Payment gateway
STRIPE_SECRET_KEY=<stripe-secret-key>
STRIPE_WEBHOOK_SECRET=<stripe-webhook-secret>
STRIPE_PUBLISHABLE_KEY=<stripe-publishable-key>
```

**Health Check**: `/health`

---

### 6. Inventory Service (`services/inventory-service`)

**Port**: 8006

```env
# Required
PORT=8006
MONGO_URI=${{MongoDB.MONGO_URL}}
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}
JWT_SECRET=<shared-secret>
ENV=production
```

**Health Check**: `/health`

---

## Railway Variable Referencing

Railway allows you to reference variables from other services in the same project:

```env
# Reference MongoDB service
MONGO_URI=${{MongoDB.MONGO_URL}}

# Reference Redis service
REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}

# Reference RabbitMQ service
RABBITMQ_URL=${{RabbitMQ.AMQP_URL}}
```

## Setting Variables in Railway

### Via Dashboard
1. Select your service
2. Go to **"Variables"** tab
3. Click **"+ New Variable"**
4. Add name and value
5. Click **"Add"**

### Via Railway CLI
```bash
# Set a single variable
railway variables set KEY=VALUE

# Set multiple variables from file
railway variables set < .env

# Set for specific service
railway variables set KEY=VALUE --service auth-service
```

### Using Shared Variables
For variables that are the same across all services (like `JWT_SECRET`):

1. Create a **"Shared Variables"** group in Railway
2. Add common variables there
3. Reference them in your services

## Generating Secure Secrets

```bash
# Generate JWT Secret (recommended: 32+ characters)
openssl rand -base64 32

# Generate a random string
openssl rand -hex 32

# Generate UUID
uuidgen
```

## Environment-Specific Configuration

### Development
```env
ENV=development
LOG_LEVEL=debug
```

### Staging
```env
ENV=staging
LOG_LEVEL=info
```

### Production
```env
ENV=production
LOG_LEVEL=warn
```

## Connection String Formats

### MongoDB
Railway format:
```
mongodb://username:password@host:port/?authSource=admin
```

If your app expects different format, you might need to transform it in your config.

### Redis
Railway provides:
```
redis://username:password@host:port
```

Your Go code should handle this or extract the host:port if needed.

### RabbitMQ (CloudAMQP)
```
amqp://username:password@host/vhost
```

## Security Best Practices

1. **Never commit secrets** to git
2. **Use Railway's built-in secrets** for sensitive data
3. **Rotate secrets regularly** (especially JWT_SECRET)
4. **Use different secrets** for dev/staging/production
5. **Limit access** to production variables
6. **Audit variable access** through Railway's logs

## Variable Validation Checklist

Before deploying, verify:

- [ ] All required variables are set
- [ ] JWT_SECRET is the same across all services
- [ ] MongoDB URI is correct and accessible
- [ ] Redis address is correct (use private URL)
- [ ] Service ports don't conflict (each service has unique port)
- [ ] External API keys are valid (Google OAuth, Stripe, etc.)
- [ ] Environment is set to "production"
- [ ] No default/example values remain

## Troubleshooting Variable Issues

### Service won't start
1. Check logs for "missing environment variable" errors
2. Verify variable names match exactly (case-sensitive)
3. Ensure no extra spaces in variable values

### Database connection fails
1. Verify `MONGO_URI` format
2. Check if MongoDB service is running in Railway
3. Test connection from Railway shell: `railway run bash`

### Redis connection fails
1. Use `REDIS_PRIVATE_URL` instead of public URL
2. Verify Redis service is in the same Railway project
3. Check if Redis requires password

### Services can't communicate
1. Ensure all services are in the same Railway project
2. Use Railway's internal URLs for service-to-service calls
3. Check firewall/security settings

## Example: Setting Up Variables for Auth Service

```bash
# Using Railway CLI
railway link  # Link to your Railway project

# Set variables for auth-service
railway variables set PORT=8001 --service auth-service
railway variables set JWT_SECRET=$(openssl rand -base64 32) --service auth-service
railway variables set JWT_EXPIRY=15m --service auth-service
railway variables set REFRESH_TOKEN_EXPIRY=168h --service auth-service
railway variables set ENV=production --service auth-service

# Reference other services
railway variables set 'MONGO_URI=${{MongoDB.MONGO_URL}}' --service auth-service
railway variables set 'REDIS_ADDR=${{Redis.REDIS_PRIVATE_URL}}' --service auth-service
```

## Additional Notes

- Railway automatically injects some variables (like `PORT`, `RAILWAY_ENVIRONMENT`)
- You can override Railway's defaults by setting your own values
- Variables are encrypted at rest
- Changes to variables trigger a redeployment
- Use Railway's UI for easier management of complex configurations

---

For more information, see:
- [RAILWAY_DEPLOYMENT.md](./RAILWAY_DEPLOYMENT.md) - Full deployment guide
- [Railway Variables Documentation](https://docs.railway.app/deploy/variables)
