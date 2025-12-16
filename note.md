- These were temperory removed from the docker-compose.yml

``` Dockerfile
  # API Gateway
  api-gateway:
    build:
      context: .
      dockerfile: ./services/api-gateway/Dockerfile
    container_name: erp-api-gateway
    restart: unless-stopped
    environment:
      PORT: 8080
      MONGO_URI: mongodb://admin:admin123@mongodb:27017
      REDIS_ADDR: redis:6379
      RABBITMQ_URL: amqp://admin:admin123@rabbitmq:5672/
      AUTH_SERVICE_URL: http://auth-service:8001
      ORG_SERVICE_URL: http://org-service:8002
      JWT_SECRET: your-super-secret-jwt-key-change-in-production
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
      - redis
      - rabbitmq
      - auth-service
      - org-service
    networks:
      - erp-network

      # Subscription Service
  subscription-service:
    build:
      context: .
      dockerfile: ./services/subscription-service/Dockerfile
    container_name: erp-subscription-service
    restart: unless-stopped
    environment:
      PORT: 8003
      MONGO_URI: mongodb://admin:admin123@mongodb:27017
      REDIS_ADDR: redis:6379
      JWT_SECRET: your-super-secret-jwt-key-change-in-production
    ports:
      - "8003:8003"
    depends_on:
      - mongodb
      - redis
    networks:
      - erp-network

      # Inventory Service
  inventory-service:
    build:
      context: .
      dockerfile: ./services/inventory-service/Dockerfile
    container_name: erp-inventory-service
    restart: unless-stopped
    environment:
      PORT: 8006
      MONGO_URI: mongodb://admin:admin123@mongodb:27017
      REDIS_ADDR: redis:6379
      JWT_SECRET: your-super-secret-jwt-key-change-in-production
    ports:
      - "8006:8006"
    depends_on:
      - mongodb
      - redis
    networks:
      - erp-network

 # License Service
  license-service:
    build:
      context: .
      dockerfile: ./services/license-service/Dockerfile
    container_name: erp-license-service
    restart: unless-stopped
    environment:
      PORT: 8004
      MONGO_URI: mongodb://admin:admin123@mongodb:27017
      REDIS_ADDR: redis:6379
      JWT_SECRET: your-super-secret-jwt-key-change-in-production
    ports:
      - "8004:8004"
    depends_on:
      - mongodb
      - redis
    networks:
      - erp-network

# AI Service (Python/FastAPI)
  ai-service:
    build:
      context: .
      dockerfile: ./services/ai-service/Dockerfile
    container_name: erp-ai-service
    restart: unless-stopped
    environment:
      PORT: 9000
      MONGO_URI: mongodb://admin:admin123@mongodb:27017
      REDIS_ADDR: redis:6379
      OPENAI_API_KEY: ${OPENAI_API_KEY}
      ANTHROPIC_API_KEY: ${ANTHROPIC_API_KEY}
    ports:
      - "9000:9000"
    depends_on:
      - mongodb
      - redis
    networks:
      - erp-network

 # Workflow Service
  workflow-service:
    build:
      context: .
      dockerfile: ./services/workflow-service/Dockerfile
    container_name: erp-workflow-service
    restart: unless-stopped
    environment:
      PORT: 8005
      MONGO_URI: mongodb://admin:admin123@mongodb:27017
      REDIS_ADDR: redis:6379
      RABBITMQ_URL: amqp://admin:admin123@rabbitmq:5672/
      AI_SERVICE_URL: http://ai-service:9000
    ports:
      - "8005:8005"
    depends_on:
      - mongodb
      - redis
      - rabbitmq
    networks:
      - erp-network


```