# Railway Quick Start Script
# This script helps you set up environment variables for Railway deployment

echo "🚂 Railway ERP System Deployment Helper"
echo "========================================"
echo ""

# Generate JWT Secret
echo "Generating JWT Secret..."
JWT_SECRET=$(openssl rand -base64 32)
echo "✓ JWT Secret generated"
echo ""

# Display generated secrets
echo "📝 Generated Secrets (SAVE THESE SECURELY):"
echo "============================================"
echo "JWT_SECRET=$JWT_SECRET"
echo ""

# Create .env.example for Railway
echo "Creating environment variable templates..."

cat > .env.railway.template << EOF
# Generated Railway Environment Variables Template
# Copy these to your Railway service settings

# ===========================================
# SHARED VARIABLES (All Services)
# ===========================================
JWT_SECRET=$JWT_SECRET
ENV=production

# ===========================================
# AUTH SERVICE (Port 8001)
# ===========================================
PORT=8001
MONGO_URI=\${{MongoDB.MONGO_URL}}
REDIS_ADDR=\${{Redis.REDIS_PRIVATE_URL}}
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
GOOGLE_CLIENT_ID=your-google-client-id-here
GOOGLE_CLIENT_SECRET=your-google-client-secret-here

# ===========================================
# ORG SERVICE (Port 8002)
# ===========================================
# PORT=8002
# MONGO_URI=\${{MongoDB.MONGO_URL}}
# REDIS_ADDR=\${{Redis.REDIS_PRIVATE_URL}}
# RABBITMQ_URL=\${{RabbitMQ.RABBITMQ_URL}}

# ===========================================
# PRODUCT SERVICE (Port 8003)
# ===========================================
# PORT=8003
# MONGO_URI=\${{MongoDB.MONGO_URL}}
# REDIS_ADDR=\${{Redis.REDIS_PRIVATE_URL}}
# RABBITMQ_URL=\${{RabbitMQ.RABBITMQ_URL}}

# ===========================================
# LICENSE SERVICE (Port 8004)
# ===========================================
# PORT=8004
# MONGO_URI=\${{MongoDB.MONGO_URL}}
# REDIS_ADDR=\${{Redis.REDIS_PRIVATE_URL}}

# ===========================================
# SUBSCRIPTION SERVICE (Port 8005)
# ===========================================
# PORT=8005
# MONGO_URI=\${{MongoDB.MONGO_URL}}
# REDIS_ADDR=\${{Redis.REDIS_PRIVATE_URL}}

# ===========================================
# INVENTORY SERVICE (Port 8006)
# ===========================================
# PORT=8006
# MONGO_URI=\${{MongoDB.MONGO_URL}}
# REDIS_ADDR=\${{Redis.REDIS_PRIVATE_URL}}

EOF

echo "✓ Environment template created: .env.railway.template"
echo ""

echo "📋 Next Steps:"
echo "=============="
echo "1. Install Railway CLI: npm install -g @railway/cli"
echo "2. Login to Railway: railway login"
echo "3. Create a new project: railway init"
echo "4. Add MongoDB: In Railway dashboard, click '+ New' → 'Database' → 'MongoDB'"
echo "5. Add Redis: In Railway dashboard, click '+ New' → 'Database' → 'Redis'"
echo "6. Deploy each service individually (see RAILWAY_DEPLOYMENT.md)"
echo ""
echo "💡 Pro Tips:"
echo "- Use Railway's shared variables for JWT_SECRET"
echo "- Reference services with \${{ServiceName.VARIABLE}}"
echo "- Set watch paths to avoid unnecessary rebuilds"
echo ""
echo "📖 For detailed instructions, see RAILWAY_DEPLOYMENT.md"
