#!/bin/bash

# Railway Monorepo Deployment Script
# This script deploys all services to Railway from a single repository

set -e

echo "🚂 Railway ERP System Deployment"
echo "=================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo -e "${RED}❌ Railway CLI not found!${NC}"
    echo "Install it with: npm install -g @railway/cli"
    exit 1
fi

echo -e "${GREEN}✓ Railway CLI found${NC}"
echo ""

# Check if logged in
if ! railway whoami &> /dev/null; then
    echo -e "${YELLOW}⚠️  Not logged in to Railway${NC}"
    echo "Logging in..."
    railway login
fi

echo -e "${GREEN}✓ Authenticated${NC}"
echo ""

# Project ID (will be set after first deployment)
RAILWAY_PROJECT_ID=${RAILWAY_PROJECT_ID:-""}

if [ -z "$RAILWAY_PROJECT_ID" ]; then
    echo -e "${YELLOW}⚠️  RAILWAY_PROJECT_ID not set${NC}"
    echo "Please set it as an environment variable or the script will prompt you"
    echo "You can find it in your Railway project settings"
    echo ""
fi

# Services to deploy
SERVICES=("auth-service" "org-service" "product-service" "inventory-service" "license-service" "subscription-service")

# Function to deploy a service
deploy_service() {
    local service_name=$1
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${YELLOW}Deploying ${service_name}...${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Change to service directory
    cd "services/${service_name}"
    
    # Deploy using Railway CLI
    if railway up --service "${service_name}"; then
        echo -e "${GREEN}✓ ${service_name} deployed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to deploy ${service_name}${NC}"
        return 1
    fi
    
    # Go back to root
    cd ../..
}

# Main deployment
echo "Select deployment option:"
echo "1) Deploy all services"
echo "2) Deploy specific service"
echo "3) Deploy changed services only"
read -p "Enter choice (1-3): " choice

case $choice in
    1)
        echo ""
        echo "Deploying all services..."
        for service in "${SERVICES[@]}"; do
            deploy_service "$service"
        done
        ;;
    2)
        echo ""
        echo "Available services:"
        for i in "${!SERVICES[@]}"; do
            echo "$((i+1))) ${SERVICES[$i]}"
        done
        read -p "Enter service number: " service_num
        service_index=$((service_num-1))
        if [ $service_index -ge 0 ] && [ $service_index -lt ${#SERVICES[@]} ]; then
            deploy_service "${SERVICES[$service_index]}"
        else
            echo -e "${RED}Invalid service number${NC}"
            exit 1
        fi
        ;;
    3)
        echo ""
        echo "Detecting changed files..."
        # Get changed files since last commit
        changed_files=$(git diff --name-only HEAD~1 HEAD 2>/dev/null || git diff --name-only --cached)
        
        if [ -z "$changed_files" ]; then
            echo -e "${YELLOW}No changes detected${NC}"
            exit 0
        fi
        
        echo "Changed files:"
        echo "$changed_files"
        echo ""
        
        # Detect which services changed
        services_to_deploy=()
        for service in "${SERVICES[@]}"; do
            if echo "$changed_files" | grep -q "services/${service}/\|shared/"; then
                services_to_deploy+=("$service")
            fi
        done
        
        if [ ${#services_to_deploy[@]} -eq 0 ]; then
            echo -e "${YELLOW}No service changes detected${NC}"
            exit 0
        fi
        
        echo "Services to deploy:"
        printf '%s\n' "${services_to_deploy[@]}"
        echo ""
        
        read -p "Continue with deployment? (y/n): " confirm
        if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
            for service in "${services_to_deploy[@]}"; do
                deploy_service "$service"
            done
        else
            echo "Deployment cancelled"
            exit 0
        fi
        ;;
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ Deployment complete!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Next steps:"
echo "1. Check Railway dashboard for deployment status"
echo "2. Verify environment variables are set"
echo "3. Test your services"
echo ""
