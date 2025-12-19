package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/services/inventory-service/config"
	"github.com/yourusername/erp-system/services/inventory-service/internal/handlers"
	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")
	db := client.Database("erp_inventory")

	// Create indexes
	if err := createIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	stockLevelRepo := repository.NewStockLevelRepository(db)
	stockMovementRepo := repository.NewStockMovementRepository(db)
	batchRepo := repository.NewBatchRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo, stockLevelRepo)
	stockService := service.NewStockService(productRepo, stockLevelRepo, stockMovementRepo, batchRepo)

	// Initialize handlers
	inventoryHandler := handlers.NewInventoryHandler(productService, stockService)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	//TODO: router.Use(middleware.RateLimitMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "inventory-service",
			"time":    time.Now(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	inventoryHandler.RegisterRoutes(api, jwtManager)

	// Start server
	port := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("Inventory Service starting on port %s", cfg.Port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Products indexes
	productIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "sku", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_product_org_sku"),
		},
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_product_org_status"),
		},
		{
			Keys:    bson.D{{Key: "barcode", Value: 1}},
			Options: options.Index().SetName("idx_product_barcode"),
		},
		{
			Keys:    bson.D{{Key: "category_id", Value: 1}},
			Options: options.Index().SetName("idx_product_category"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_product_deleted"),
		},
	}
	if _, err := db.Collection("products").Indexes().CreateMany(ctx, productIndexes); err != nil {
		return fmt.Errorf("failed to create product indexes: %w", err)
	}

	// Stock levels indexes
	stockIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "product_id", Value: 1}, {Key: "location_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_stock_org_product_location"),
		},
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "location_id", Value: 1}},
			Options: options.Index().SetName("idx_stock_product_location"),
		},
		{
			Keys:    bson.D{{Key: "location_id", Value: 1}},
			Options: options.Index().SetName("idx_stock_location"),
		},
	}
	if _, err := db.Collection("stock_levels").Indexes().CreateMany(ctx, stockIndexes); err != nil {
		return fmt.Errorf("failed to create stock indexes: %w", err)
	}

	// Stock movements indexes
	movementIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "movement_date", Value: -1}},
			Options: options.Index().SetName("idx_movement_org_date"),
		},
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "movement_date", Value: -1}},
			Options: options.Index().SetName("idx_movement_product_date"),
		},
		{
			Keys:    bson.D{{Key: "reference_type", Value: 1}, {Key: "reference_id", Value: 1}},
			Options: options.Index().SetName("idx_movement_reference"),
		},
	}
	if _, err := db.Collection("stock_movements").Indexes().CreateMany(ctx, movementIndexes); err != nil {
		return fmt.Errorf("failed to create movement indexes: %w", err)
	}

	// Batches indexes
	batchIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "batch_number", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_batch_org_number"),
		},
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_batch_product_active"),
		},
		{
			Keys:    bson.D{{Key: "expiry_date", Value: 1}},
			Options: options.Index().SetName("idx_batch_expiry"),
		},
	}
	if _, err := db.Collection("batches").Indexes().CreateMany(ctx, batchIndexes); err != nil {
		return fmt.Errorf("failed to create batch indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
