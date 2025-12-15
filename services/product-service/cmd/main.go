package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/product-service/config"
	"github.com/yourusername/erp-system/services/product-service/internal/handlers"
	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"github.com/yourusername/erp-system/services/product-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yourusername/erp-system/shared/database"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)


func main() {

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
		// Connect to MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI, "erp_db")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Connect to Redis
	redisClient, err := database.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create indexes
	if err := createIndexes(mongoDB.Database); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 15*time.Minute, 168*time.Hour)


	// Initialize repositories
	productRepo := repository.NewProductRepository(mongoDB.Database)
	
	// Initialize services
	productService := service.NewProductService(productRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(redisClient.Client, 100, time.Minute)
	router.Use(rateLimiter.Middleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, 200, gin.H{"status": "healthy"}, "Service is healthy")
	})

	// API routes
	v1 := router.Group("/api/v1")
	productHandler.RegisterProductRoutes(v1, jwtManager)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Product Service started on port %s", cfg.Port)


		// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func createIndexes(db *mongo.Database) error {
	//TODO: Create necessary indexes for product collections
	return nil
}