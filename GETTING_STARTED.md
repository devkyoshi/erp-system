# Getting Started with ERP System Backend

## What's Been Implemented

### ✅ Core Infrastructure
1. **Project Structure**: Microservices architecture with shared libraries
2. **Docker Compose**: Complete local development environment
3. **Shared Libraries**:
   - MongoDB and Redis connection managers
   - JWT authentication utilities
   - Password hashing
   - Response formatters
   - Middleware (Auth, CORS, Rate Limiting)
   - Data models for all entities

### ✅ Authentication Service (COMPLETE)
Full-featured auth service with:
- User registration with email/password
- User login with JWT tokens
- Google OAuth 2.0 integration
- Access token + Refresh token mechanism
- Token refresh endpoint
- User profile endpoint
- Session management
- Password hashing with bcrypt
- Rate limiting
- CORS support

**Endpoints**:
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login with email/password
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout (invalidate session)
- `GET /api/v1/auth/google` - Get Google OAuth URL
- `GET /api/v1/auth/google/callback` - Google OAuth callback
- `GET /api/v1/auth/me` - Get current user profile (protected)

## Quick Start

### 1. Prerequisites
```bash
# Install Go 1.21+
go version

# Install Docker & Docker Compose
docker --version
docker-compose --version

# Install Python 3.11+ (for AI services later)
python --version
```

### 2. Clone and Setup
```bash
cd /home/user/erp-system-backend

# Create .env file
cp .env.example .env

# Edit .env and add your credentials:
# - GOOGLE_CLIENT_ID
# - GOOGLE_CLIENT_SECRET
# - OPENAI_API_KEY (for later)
# - ANTHROPIC_API_KEY (for later)
```

### 3. Start Infrastructure
```bash
# Start MongoDB, Redis, RabbitMQ
docker-compose up -d mongodb redis rabbitmq

# Wait for services to be ready (30 seconds)
sleep 30
```

### 4. Install Dependencies
```bash
# Install shared dependencies
cd shared && go mod tidy && cd ..

# Install auth service dependencies
cd services/auth-service && go mod download && cd ../..
```

### 5. Run Auth Service
```bash
cd services/auth-service
go run cmd/main.go
```

The service will start on `http://localhost:8001`

### 6. Test the Auth Service

#### Register a new user:
```bash
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePass123",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

#### Login:
```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePass123"
  }'
```

You'll receive:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "user": { ... }
  },
  "message": "Login successful",
  "timestamp": "2025-12-07T..."
}
```

#### Get current user (use access_token from login):
```bash
curl -X GET http://localhost:8001/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Next Implementation Steps

### Phase 1: Organization Service (Next)
Implement the organization hierarchy:
- Organizations
- Companies (under organizations)
- Locations (under companies)
- User-location mappings

### Phase 2: License & Subscription Services
- Application definitions
- License management (user-based, device-based, usage-based)
- Subscription plans
- Feature flags

### Phase 3: Workflow Engine
- Workflow definitions (JSON DSL)
- Workflow execution engine
- State management
- Human tasks
- AI agent integration points

### Phase 4: Dynamic Form Builder
- JSON Schema-based forms
- Form validation
- Conditional logic
- AI-assisted form generation

### Phase 5: Core ERP Modules
- Product & Inventory Management
- Purchase Order & GRN
- Supplier Management
- CRM
- Sales & Invoices
- Bill Management

### Phase 6: AI Service (Python/FastAPI)
- Natural language query interface
- RAG implementation for ERP data
- Workflow suggestions
- Anomaly detection
- Predictive analytics

### Phase 7: API Gateway
- Request routing
- Authentication gateway
- Rate limiting
- API documentation (Swagger)

## Development Commands

### Using Makefile
```bash
# Install all dependencies
make install

# Run all services
make dev

# Build all services
make build

# Run tests
make test

# Clean build artifacts
make clean

# Database operations
make migrate
make seed

# Code quality
make lint
make format
```

### Docker Commands
```bash
# Start all infrastructure
make docker-up

# Stop all services
make docker-down

# Clean everything (including volumes)
make docker-clean
```

## Project Structure

```
erp-system-backend/
├── services/
│   └── auth-service/          ✅ COMPLETE
│       ├── cmd/main.go
│       ├── config/
│       ├── internal/
│       │   ├── handlers/
│       │   ├── repository/
│       │   └── service/
│       ├── Dockerfile
│       └── go.mod
├── shared/                    ✅ COMPLETE
│   ├── database/
│   ├── models/
│   ├── middleware/
│   ├── utils/
│   └── go.mod
├── ARCHITECTURE.md           ✅ Complete design doc
├── docker-compose.yml        ✅ Full stack
├── Makefile                  ✅ Build automation
└── README.md                 ✅ Documentation
```

## Database Schema

### Collections Created
- `users` - User accounts
- `sessions` - Active sessions with refresh tokens
- `roles` - User roles (planned)
- `permissions` - Permissions (planned)
- `organizations` - Organizations (next)
- `companies` - Companies (next)
- `locations` - Locations (next)

### Indexes Created
Auth service creates the following indexes:
- `users.email` (unique)
- `users.organization_id`
- `users.google_id`
- `sessions.refresh_token`
- `sessions.user_id`
- `sessions.expires_at`

## Security Features

✅ **Implemented**:
- JWT-based authentication (RS256 recommended for production)
- Password hashing with bcrypt
- Access tokens (15 min expiry)
- Refresh tokens (7 days expiry)
- Session management
- Rate limiting (100 req/min per user)
- CORS protection
- Input validation

🔄 **Planned**:
- MFA (Multi-Factor Authentication)
- API key authentication for AI agents
- Role-Based Access Control (RBAC)
- Permission-based authorization
- Audit logging
- Field-level encryption for PII

## Monitoring & Logging

🔄 **Planned Integration**:
- Prometheus metrics
- Grafana dashboards
- ELK stack for logging
- Health check endpoints (implemented)
- Performance monitoring

## Testing Strategy

```bash
# Unit tests
go test ./services/auth-service/...

# Integration tests
go test -tags=integration ./services/auth-service/...

# Load tests (k6)
k6 run scripts/load-tests/auth.js
```

## Common Issues & Troubleshooting

### MongoDB Connection Failed
```bash
# Check if MongoDB is running
docker ps | grep mongo

# Check logs
docker logs erp-mongodb

# Restart MongoDB
docker-compose restart mongodb
```

### Redis Connection Failed
```bash
# Check if Redis is running
docker ps | grep redis

# Test connection
docker exec -it erp-redis redis-cli ping
```

### Port Already in Use
```bash
# Find process using port 8001
lsof -i :8001

# Kill process
kill -9 <PID>
```

## Production Deployment

### Kubernetes Deployment (Planned)
```bash
# Deploy to k8s
kubectl apply -f infrastructure/kubernetes/

# Check status
kubectl get pods -n erp-system
```

### Environment Variables for Production
- Change `JWT_SECRET` to a strong random string
- Use managed MongoDB (Atlas) or self-hosted cluster
- Use managed Redis (ElastiCache, Redis Cloud)
- Enable TLS/SSL
- Configure proper CORS origins
- Set up monitoring and alerting

## API Documentation

Once API Gateway is implemented, Swagger docs will be available at:
- Development: `http://localhost:8080/swagger`
- Production: `https://api.yourdomain.com/swagger`

## Contributing

### Code Standards
- Follow Go best practices
- Write unit tests for new features
- Update documentation
- Use conventional commits
- Run `make lint` before committing

### Git Workflow
```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make changes and commit
git add .
git commit -m "feat: add new feature"

# Push to remote
git push -u origin feature/your-feature-name
```

## Support

For questions or issues:
1. Check this documentation
2. Review ARCHITECTURE.md
3. Check existing GitHub issues
4. Create a new issue with detailed description

## License

Proprietary - All rights reserved

---

## What's Next?

Would you like me to:
1. **Continue building** the Organization Service?
2. **Add more features** to the Auth Service (MFA, email verification)?
3. **Create the API Gateway** to route requests?
4. **Jump to the Workflow Engine** (most advanced feature)?
5. **Start the AI Service** in Python?
6. **Build specific ERP modules** (Product, Inventory, etc.)?

Let me know your priority and I'll continue implementing!
