# Gemini Code Assistant Context

This document provides a comprehensive overview of the ERP system backend project for the Gemini Code Assistant.

## Project Overview

This project is a multi-tenant, AI-ready, workflow-driven ERP system. It is built with a microservices architecture using Golang for backend services and Python for AI-related tasks.

### Key Technologies

- **Backend:** Golang (Gin framework)
- **AI Services:** Python (FastAPI)
- **Database:** MongoDB
- **Cache:** Redis
- **Message Queue:** RabbitMQ
- **Containerization:** Docker and Docker Compose

### Architecture

The system is composed of multiple microservices that communicate with each other through a message queue (RabbitMQ). A full breakdown of the architecture can be found in `ARCHITECTURE.md`.

## Building and Running the Project

The project uses `make` and `docker-compose` to streamline the development process.

### First-Time Setup

1.  **Start Infrastructure:**
    ```bash
    docker-compose up -d mongodb redis rabbitmq
    ```

2.  **Install Dependencies:**
    ```bash
    make install
    ```
3. **Seed the database**
   ```bash
   make seed-all
   ```

### Running the Application

To start all the services, run:

```bash
make dev
```

This command is an alias for `docker-compose up`.

### Running Tests

To run the test suite for all services, use:

```bash
make test
```

### Other Useful Commands

-   `make build`: Build all services.
-   `make clean`: Clean build artifacts.
-   `make docker-down`: Stop all Docker services.
-   `make swagger-serve`: Serve the Swagger UI for API documentation.

## Development Conventions

-   **Microservices:** The project is structured as a collection of microservices, each residing in its own directory under `services/`.
-   **Shared Code:** Common code shared across services is located in the `shared/` directory.
-   **Dependency Management:** Go modules are used for managing dependencies. Each service has its own `go.mod` file.
-   **Makefile:** A `Makefile` is provided to automate common development tasks.
-   **Containerization:** All services are containerized using Docker. `docker-compose.yml` defines the services and their configurations.

This information should provide a solid foundation for any requests regarding the project. For more in-depth information, please refer to the `README.md` and `ARCHITECTURE.md` files.
