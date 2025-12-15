# Advanced ERP System Backend

An AI-ready, workflow-driven, multi-tenant ERP system built with Golang, MongoDB, and Python.

## Features

- **Multi-Tenant Architecture**: Organization → Company → Location hierarchy
- **Workflow Engine**: Customizable business workflows
- **Dynamic Forms**: JSON Schema-based form builder
- **AI Integration**: Natural language queries, automation, and analytics
- **License Management**: User-based, device-based, and usage-based licensing
- **Subscription Plans**: Flexible plan management
- **Complete ERP Modules**: Products, Inventory, PO, GRN, Suppliers, CRM, Sales, Bills

## Tech Stack

- **Backend**: Golang (Gin framework)
- **Database**: MongoDB
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **AI Services**: Python (FastAPI)
- **Containerization**: Docker + Docker Compose

## Quick Start

### Prerequisites
- Go 1.21+
- Python 3.11+
- Docker & Docker Compose
- MongoDB 7.0+
- Redis 7.0+

### Installation

1. Clone the repository
```bash
git clone <repo-url>
cd erp-system-backend
```

2. Start infrastructure services
```bash
docker-compose up -d mongodb redis rabbitmq
```

3. Install dependencies
```bash
make install
```

4. Run migrations
```bash
make migrate
```

5. Start services
```bash
make dev
```

## Project Structure

```
erp-system-backend/
├── services/          # Microservices
├── shared/            # Shared libraries
├── infrastructure/    # Docker, K8s configs
├── scripts/           # Utility scripts
└── docs/              # Documentation
```

## Services

- **api-gateway**: API Gateway with routing and auth
- **auth-service**: Authentication & authorization
- **org-service**: Organization hierarchy management
- **license-service**: License management
- **subscription-service**: Subscription plans
- **workflow-service**: Workflow engine
- **form-service**: Dynamic form builder
- **product-service**: Product & inventory management
- **purchase-service**: PO & GRN
- **supplier-service**: Supplier management
- **crm-service**: Customer relationship management
- **sales-service**: Sales & invoices
- **bill-service**: Bill management
- **notification-service**: Multi-channel notifications
- **ai-service**: AI agent (Python)

## API Documentation

API documentation is available at:
- Development: http://localhost:8080/swagger
- Production: https://api.yourdomain.com/swagger

## License

Proprietary
