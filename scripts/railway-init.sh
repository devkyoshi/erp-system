#!/bin/bash

# Railway Setup Script - One-time setup for Railway deployment

set -e

echo "🚂 Railway Project Setup"
echo "========================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check Railway CLI
if ! command -v railway &> /dev/null; then
    echo -e "${RED}❌ Railway CLI not installed${NC}"
    echo "Installing Railway CLI..."
    npm install -g @railway/cli
fi

echo -e "${GREEN}✓ Railway CLI ready${NC}"
echo ""

# Login
echo "Step 1: Login to Railway"
echo "-------------------------"
railway login
echo ""

# Create or link project
echo "Step 2: Project Setup"
echo "---------------------"
echo "Choose an option:"
echo "1) Create new Railway project"
echo "2) Link to existing project"
read -p "Enter choice (1-2): " project_choice

case $project_choice in
    1)
        echo ""
        read -p "Enter project name: " project_name
        railway init --name "$project_name"
        ;;
    2)
        echo ""
        echo "Go to your Railway dashboard and copy your project ID"
        read -p "Enter project ID: " project_id
        railway link "$project_id"
        ;;
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}✓ Project linked${NC}"
echo ""

# Get project info
PROJECT_ID=$(railway status | grep "Project ID" | awk '{print $3}' || echo "")
echo "Project ID: $PROJECT_ID"
echo ""

# Generate secrets
echo "Step 3: Generate Secrets"
echo "------------------------"
JWT_SECRET=$(openssl rand -base64 32)
echo -e "${GREEN}✓ JWT Secret generated${NC}"
echo ""

# Save to .env file
cat > .env.railway << EOF
# Railway Project Configuration
RAILWAY_PROJECT_ID=$PROJECT_ID

# Generated Secrets (KEEP SECURE!)
JWT_SECRET=$JWT_SECRET

# Service Ports
AUTH_SERVICE_PORT=8001
ORG_SERVICE_PORT=8002
PRODUCT_SERVICE_PORT=8003
LICENSE_SERVICE_PORT=8004
SUBSCRIPTION_SERVICE_PORT=8005
INVENTORY_SERVICE_PORT=8006
EOF

echo -e "${GREEN}✓ Configuration saved to .env.railway${NC}"
echo -e "${YELLOW}⚠️  Keep this file secure and do not commit it!${NC}"
echo ""

# Create services
echo "Step 4: Create Services"
echo "-----------------------"
echo ""
echo "Now you need to create each service in Railway:"
echo ""
echo "For each service (auth, org, product, inventory, license, subscription):"
echo "1. Go to your Railway dashboard"
echo "2. Click '+ New' → 'Empty Service'"
echo "3. Name it (e.g., 'auth-service')"
echo "4. Connect your GitHub repo"
echo "5. Set Root Directory to: / (repository root)"
echo "6. The nixpacks.toml will handle the rest"
echo ""

read -p "Press Enter when you've created all services in Railway..."

echo ""
echo "Step 5: Set Environment Variables"
echo "----------------------------------"
echo ""
echo "Setting shared environment variables..."

# Set JWT_SECRET for all services
SERVICES=("auth-service" "org-service" "product-service" "inventory-service" "license-service" "subscription-service")

for service in "${SERVICES[@]}"; do
    echo "Setting variables for $service..."
    
    # This would require the service to exist first
    # railway variables set JWT_SECRET="$JWT_SECRET" --service "$service" 2>/dev/null || echo "  → Skip (set manually)"
    echo "  → Set JWT_SECRET manually in Railway dashboard"
done

echo ""
echo "You'll need to set these variables manually in Railway dashboard:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "For ALL services:"
echo "  JWT_SECRET=$JWT_SECRET"
echo "  ENV=production"
echo ""
echo "Service-specific (see ENVIRONMENT_VARIABLES.md for details):"
echo "  - PORT (different for each service)"
echo "  - MONGO_URI=\${{MongoDB.MONGO_URL}}"
echo "  - REDIS_ADDR=\${{Redis.REDIS_PRIVATE_URL}}"
echo ""

# Add databases
echo "Step 6: Add Databases"
echo "---------------------"
echo ""
echo "In your Railway project:"
echo "1. Click '+ New' → 'Database' → 'Add MongoDB'"
echo "2. Click '+ New' → 'Database' → 'Add Redis'"
echo ""
echo "Railway will automatically create connection strings"
echo ""

read -p "Press Enter when you've added MongoDB and Redis..."

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ Setup Complete!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 Important files created:"
echo "  - .env.railway (contains your secrets)"
echo ""
echo "🚀 Next steps:"
echo "  1. Set environment variables in Railway dashboard"
echo "  2. Deploy services: bash scripts/deploy-railway.sh"
echo "  3. Monitor deployments in Railway dashboard"
echo ""
echo "📖 Documentation:"
echo "  - RAILWAY_DEPLOYMENT.md - Detailed guide"
echo "  - ENVIRONMENT_VARIABLES.md - All env vars"
echo ""
