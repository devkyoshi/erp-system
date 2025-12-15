package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type InventoryHandler struct {
	productService *service.ProductService
	stockService   *service.StockService
}

func NewInventoryHandler(productService *service.ProductService, stockService *service.StockService) *InventoryHandler {
	return &InventoryHandler{
		productService: productService,
		stockService:   stockService,
	}
}

// Product Handlers

func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	product, err := h.productService.CreateProduct(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, product, "Product created successfully")
}

func (h *InventoryHandler) GetProduct(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, product, "Product retrieved successfully")
}

func (h *InventoryHandler) ListProducts(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		filters["category_id"] = categoryID
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	products, err := h.productService.ListProducts(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, products, "Products retrieved successfully")
}

func (h *InventoryHandler) UpdateProduct(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	updatedBy, _ := primitive.ObjectIDFromHex(userID)

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, product, "Product updated successfully")
}

func (h *InventoryHandler) DeleteProduct(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.productService.DeleteProduct(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Product deleted successfully")
}

// Stock Handlers

func (h *InventoryHandler) CreateStockMovement(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.StockMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	movement, err := h.stockService.CreateStockMovement(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, movement, "Stock movement created successfully")
}

func (h *InventoryHandler) GetStockLevels(c *gin.Context) {
	productIDStr := c.Query("product_id")
	locationIDStr := c.Query("location_id")

	if productIDStr == "" || locationIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "product_id and location_id required", nil)
		return
	}

	productID, _ := primitive.ObjectIDFromHex(productIDStr)
	locationID, _ := primitive.ObjectIDFromHex(locationIDStr)

	stock, err := h.stockService.GetStockLevels(c.Request.Context(), productID, locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stock, "Stock levels retrieved successfully")
}

func (h *InventoryHandler) GetStockByLocation(c *gin.Context) {
	locationIDStr := c.Param("location_id")
	locationID, err := primitive.ObjectIDFromHex(locationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location ID", nil)
		return
	}

	stocks, err := h.stockService.GetStockByLocation(c.Request.Context(), locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stocks, "Stock retrieved successfully")
}

func (h *InventoryHandler) GetStockMovements(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}


	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	movements, err := h.stockService.GetStockMovements(c.Request.Context(), productID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movements, "Stock movements retrieved successfully")
}

func (h *InventoryHandler) GetBatches(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	locationIDStr := c.Query("location_id")
	if locationIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "location_id required", nil)
		return
	}

	locationID, _ := primitive.ObjectIDFromHex(locationIDStr)
	activeOnly := c.Query("active") == "true"

	batches, err := h.stockService.GetBatches(c.Request.Context(), productID, locationID, activeOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, batches, "Batches retrieved successfully")
}

// Register all routes
func (h *InventoryHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Product routes
	protected.POST("/organizations/:org_id/products", h.CreateProduct)
	protected.GET("/organizations/:org_id/products", h.ListProducts)
	protected.GET("/products/:id", h.GetProduct)
	protected.PUT("/products/:id", h.UpdateProduct)
	protected.DELETE("/products/:id", h.DeleteProduct)

	// Stock routes
	protected.POST("/organizations/:org_id/stock-movements", h.CreateStockMovement)
	protected.GET("/stock-levels", h.GetStockLevels)
	protected.GET("/locations/:location_id/stock", h.GetStockByLocation)
	protected.GET("/products/:product_id/movements", h.GetStockMovements)
	protected.GET("/products/:product_id/batches", h.GetBatches)
}
