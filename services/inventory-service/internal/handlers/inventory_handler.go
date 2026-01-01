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
	stockLevelService   *service.StockLevelService
	stockService        *service.StockService
	batchService        *service.BatchService
	adjustmentService   *service.StockAdjustmentService
	countService        *service.InventoryCountService
	serialNumberService *service.SerialNumberService
}

func NewInventoryHandler(
	stockLevelService *service.StockLevelService,
	stockService *service.StockService,
	batchService *service.BatchService,
	adjustmentService *service.StockAdjustmentService,
	countService *service.InventoryCountService,
	serialNumberService *service.SerialNumberService,
) *InventoryHandler {
	return &InventoryHandler{
		stockLevelService:   stockLevelService,
		stockService:        stockService,
		batchService:        batchService,
		adjustmentService:   adjustmentService,
		countService:        countService,
		serialNumberService: serialNumberService,
	}
}

// ==================== Stock Level Handlers ====================

func (h *InventoryHandler) GetStockLevel(c *gin.Context) {
	productIDStr := c.Query("product_id")
	locationIDStr := c.Query("location_id")

	if productIDStr == "" || locationIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "product_id and location_id required", nil)
		return
	}

	productID, _ := primitive.ObjectIDFromHex(productIDStr)
	locationID, _ := primitive.ObjectIDFromHex(locationIDStr)

	stock, err := h.stockLevelService.GetStockLevel(c.Request.Context(), productID, locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stock, "Stock level retrieved successfully")
}

func (h *InventoryHandler) GetStockByProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	stocks, err := h.stockLevelService.GetStockByProduct(c.Request.Context(), productID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stocks, "Stock retrieved successfully")
}

func (h *InventoryHandler) GetStockByLocation(c *gin.Context) {
	locationIDStr := c.Param("location_id")
	locationID, err := primitive.ObjectIDFromHex(locationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location ID", nil)
		return
	}

	stocks, err := h.stockLevelService.GetStockByLocation(c.Request.Context(), locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stocks, "Stock retrieved successfully")
}

// ==================== Stock Movement Handlers ====================

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

// ==================== Batch Handlers ====================

func (h *InventoryHandler) CreateBatch(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	batch, err := h.batchService.CreateBatch(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, batch, "Batch created successfully")
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

	batches, err := h.batchService.GetBatchesByProduct(c.Request.Context(), productID, locationID, activeOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, batches, "Batches retrieved successfully")
}

func (h *InventoryHandler) GetBatch(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid batch ID", nil)
		return
	}

	batch, err := h.batchService.GetBatch(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, batch, "Batch retrieved successfully")
}

func (h *InventoryHandler) UpdateBatchQuantity(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid batch ID", nil)
		return
	}

	var req struct {
		Delta float64 `json:"delta" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.batchService.UpdateBatchQuantity(c.Request.Context(), id, req.Delta); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Batch quantity updated successfully")
}

// ==================== Stock Adjustment Handlers ====================

func (h *InventoryHandler) CreateStockAdjustment(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateStockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	adjustment, err := h.adjustmentService.CreateAdjustment(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, adjustment, "Stock adjustment created successfully")
}

func (h *InventoryHandler) GetStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	adjustment, err := h.adjustmentService.GetAdjustment(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adjustment, "Adjustment retrieved successfully")
}

func (h *InventoryHandler) UpdateStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	var req service.CreateStockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	updatedBy, _ := primitive.ObjectIDFromHex(userID)

	adjustment, err := h.adjustmentService.UpdateAdjustment(c.Request.Context(), id, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adjustment, "Adjustment updated successfully")
}

func (h *InventoryHandler) ListStockAdjustments(c *gin.Context) {
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
	if locationID := c.Query("location_id"); locationID != "" {
		filters["location_id"] = locationID
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	adjustments, err := h.adjustmentService.ListAdjustments(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adjustments, "Adjustments retrieved successfully")
}

func (h *InventoryHandler) ApproveStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	approvedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.adjustmentService.ApproveAdjustment(c.Request.Context(), id, approvedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "APPROVE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Adjustment approved successfully")
}

func (h *InventoryHandler) RejectStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	rejectedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.adjustmentService.RejectAdjustment(c.Request.Context(), id, rejectedBy, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "REJECT_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Adjustment rejected successfully")
}

func (h *InventoryHandler) DeleteStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.adjustmentService.DeleteAdjustment(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Adjustment deleted successfully")
}

// ==================== Inventory Count Handlers ====================

func (h *InventoryHandler) CreateInventoryCount(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateInventoryCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	count, err := h.countService.CreateCount(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, count, "Inventory count created successfully")
}

func (h *InventoryHandler) GetInventoryCount(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	count, err := h.countService.GetCount(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, count, "Count retrieved successfully")
}

func (h *InventoryHandler) ListInventoryCounts(c *gin.Context) {
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
	if locationID := c.Query("location_id"); locationID != "" {
		filters["location_id"] = locationID
	}
	if countType := c.Query("count_type"); countType != "" {
		filters["count_type"] = countType
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	counts, err := h.countService.ListCounts(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, counts, "Counts retrieved successfully")
}

func (h *InventoryHandler) UpdateCountItem(c *gin.Context) {
	countID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	var req service.UpdateCountItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	countedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.countService.UpdateCountItem(c.Request.Context(), countID, req, countedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count item updated successfully")
}

func (h *InventoryHandler) CompleteInventoryCount(c *gin.Context) {
	countID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	var req struct {
		CreateAdjustment bool `json:"create_adjustment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.CreateAdjustment = false
	}

	userID := middleware.GetUserID(c)
	completedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.countService.CompleteCount(c.Request.Context(), countID, completedBy, req.CreateAdjustment); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "COMPLETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count completed successfully")
}

func (h *InventoryHandler) CancelInventoryCount(c *gin.Context) {
	countID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	if err := h.countService.CancelCount(c.Request.Context(), countID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CANCEL_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count cancelled successfully")
}

func (h *InventoryHandler) DeleteInventoryCount(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.countService.DeleteCount(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count deleted successfully")
}

// ==================== Serial Number Handlers ====================

func (h *InventoryHandler) CreateSerialNumber(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateSerialNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	serialNumber, err := h.serialNumberService.CreateSerialNumber(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, serialNumber, "Serial number created successfully")
}

func (h *InventoryHandler) GetSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	serialNumber, err := h.serialNumberService.GetSerialNumber(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, serialNumber, "Serial number retrieved successfully")
}

func (h *InventoryHandler) UpdateSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.serialNumberService.UpdateSerialNumber(c.Request.Context(), id, updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number updated successfully")
}

func (h *InventoryHandler) AllocateSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	var req struct {
		CustomerID   string `json:"customer_id" binding:"required"`
		SalesOrderID string `json:"sales_order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	customerID, _ := primitive.ObjectIDFromHex(req.CustomerID)
	salesOrderID, _ := primitive.ObjectIDFromHex(req.SalesOrderID)

	if err := h.serialNumberService.AllocateSerialNumber(c.Request.Context(), id, customerID, salesOrderID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ALLOCATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number allocated successfully")
}

func (h *InventoryHandler) MarkSerialAsSold(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	if err := h.serialNumberService.MarkAsSold(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number marked as sold successfully")
}

func (h *InventoryHandler) DeleteSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.serialNumberService.DeleteSerialNumber(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number deleted successfully")
}

func (h *InventoryHandler) GetSerialNumberBySerial(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	serialNo := c.Param("serial_no")
	if serialNo == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "serial_no required", nil)
		return
	}

	serialNumber, err := h.serialNumberService.GetSerialNumberBySerial(c.Request.Context(), orgID, serialNo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, serialNumber, "Serial number retrieved successfully")
}

// ==================== Stock Movement Additional Handlers ====================

func (h *InventoryHandler) GetStockMovement(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid movement ID", nil)
		return
	}

	movement, err := h.stockService.GetStockMovement(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movement, "Stock movement retrieved successfully")
}

func (h *InventoryHandler) ListStockMovements(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if movementType := c.Query("movement_type"); movementType != "" {
		filters["movement_type"] = movementType
	}
	if locationID := c.Query("location_id"); locationID != "" {
		filters["location_id"] = locationID
	}
	if productID := c.Query("product_id"); productID != "" {
		filters["product_id"] = productID
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	movements, err := h.stockService.ListStockMovements(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movements, "Stock movements retrieved successfully")
}

func (h *InventoryHandler) GetStockMovementsByLocation(c *gin.Context) {
	locationIDStr := c.Param("location_id")
	locationID, err := primitive.ObjectIDFromHex(locationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location ID", nil)
		return
	}

	movements, err := h.stockService.GetStockMovementsByLocation(c.Request.Context(), locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movements, "Stock movements retrieved successfully")
}

// ==================== Stock Level Additional Handlers ====================

func (h *InventoryHandler) AllocateStock(c *gin.Context) {
	var req struct {
		ProductID  string  `json:"product_id" binding:"required"`
		LocationID string  `json:"location_id" binding:"required"`
		Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	productID, _ := primitive.ObjectIDFromHex(req.ProductID)
	locationID, _ := primitive.ObjectIDFromHex(req.LocationID)

	if err := h.stockLevelService.AllocateStock(c.Request.Context(), productID, locationID, req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ALLOCATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Stock allocated successfully")
}

func (h *InventoryHandler) ReleaseStock(c *gin.Context) {
	var req struct {
		ProductID  string  `json:"product_id" binding:"required"`
		LocationID string  `json:"location_id" binding:"required"`
		Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	productID, _ := primitive.ObjectIDFromHex(req.ProductID)
	locationID, _ := primitive.ObjectIDFromHex(req.LocationID)

	if err := h.stockLevelService.ReleaseStock(c.Request.Context(), productID, locationID, req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "RELEASE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Stock released successfully")
}

func (h *InventoryHandler) ListSerialNumbers(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if locationID := c.Query("location_id"); locationID != "" {
		filters["location_id"] = locationID
	}
	if available := c.Query("available"); available == "true" {
		filters["is_available"] = true
	} else if available == "false" {
		filters["is_available"] = false
	}

	serialNumbers, err := h.serialNumberService.ListSerialNumbersByProduct(c.Request.Context(), productID, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, serialNumbers, "Serial numbers retrieved successfully")
}

// Register all routes
func (h *InventoryHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Stock level routes
	protected.GET("/stock-levels", h.GetStockLevel)
	protected.GET("/products/:product_id/stock", h.GetStockByProduct)
	protected.GET("/locations/:location_id/stock", h.GetStockByLocation)
	protected.POST("/stock/allocate", h.AllocateStock)
	protected.POST("/stock/release", h.ReleaseStock)

	// Stock movement routes
	protected.POST("/organizations/:org_id/stock-movements", h.CreateStockMovement)
	protected.GET("/organizations/:org_id/stock-movements", h.ListStockMovements)
	protected.GET("/stock-movements/:id", h.GetStockMovement)
	protected.GET("/products/:product_id/movements", h.GetStockMovements)
	protected.GET("/locations/:location_id/movements", h.GetStockMovementsByLocation)

	// Batch routes
	protected.POST("/organizations/:org_id/batches", h.CreateBatch)
	protected.GET("/batches/:id", h.GetBatch)
	protected.PUT("/batches/:id/quantity", h.UpdateBatchQuantity)
	protected.GET("/products/:product_id/batches", h.GetBatches)

	// Stock adjustment routes
	protected.POST("/organizations/:org_id/stock-adjustments", h.CreateStockAdjustment)
	protected.GET("/organizations/:org_id/stock-adjustments", h.ListStockAdjustments)
	protected.GET("/stock-adjustments/:id", h.GetStockAdjustment)
	protected.PUT("/stock-adjustments/:id", h.UpdateStockAdjustment)
	protected.POST("/stock-adjustments/:id/approve", h.ApproveStockAdjustment)
	protected.POST("/stock-adjustments/:id/reject", h.RejectStockAdjustment)
	protected.DELETE("/stock-adjustments/:id", h.DeleteStockAdjustment)

	// Inventory count routes
	protected.POST("/organizations/:org_id/inventory-counts", h.CreateInventoryCount)
	protected.GET("/organizations/:org_id/inventory-counts", h.ListInventoryCounts)
	protected.GET("/inventory-counts/:id", h.GetInventoryCount)
	protected.POST("/inventory-counts/:id/items", h.UpdateCountItem)
	protected.POST("/inventory-counts/:id/complete", h.CompleteInventoryCount)
	protected.POST("/inventory-counts/:id/cancel", h.CancelInventoryCount)
	protected.DELETE("/inventory-counts/:id", h.DeleteInventoryCount)

	// Serial number routes
	protected.POST("/organizations/:org_id/serial-numbers", h.CreateSerialNumber)
	protected.GET("/serial-numbers/:id", h.GetSerialNumber)
	protected.GET("/organizations/:org_id/serial-numbers/:serial_no", h.GetSerialNumberBySerial)
	protected.GET("/products/:product_id/serial-numbers", h.ListSerialNumbers)
	protected.PUT("/serial-numbers/:id", h.UpdateSerialNumber)
	protected.POST("/serial-numbers/:id/allocate", h.AllocateSerialNumber)
	protected.POST("/serial-numbers/:id/sold", h.MarkSerialAsSold)
	protected.DELETE("/serial-numbers/:id", h.DeleteSerialNumber)
}
