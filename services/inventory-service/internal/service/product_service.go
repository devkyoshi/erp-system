package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
)

type ProductService struct {
	productRepo *repository.ProductRepository
	stockRepo   *repository.StockLevelRepository
}

func NewProductService(productRepo *repository.ProductRepository, stockRepo *repository.StockLevelRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		stockRepo:   stockRepo,
	}
}

type CreateProductRequest struct {
	SKU                 string                       `json:"sku" binding:"required"`
	Barcode             string                       `json:"barcode"`
	Name                string                       `json:"name" binding:"required"`
	Description         string                       `json:"description"`
	Type                models.ProductType           `json:"type"`
	CategoryID          string                       `json:"category_id"`
	TrackBatches        bool                         `json:"track_batches"`
	TrackSerialNumbers  bool                         `json:"track_serial_numbers"`
	ValuationMethod     models.StockValuationMethod  `json:"valuation_method"`
	BaseUOM             string                       `json:"base_uom" binding:"required"`
	Weight              float64                      `json:"weight"`
	CostPrice           float64                      `json:"cost_price"`
	SellingPrice        float64                      `json:"selling_price"`
	ReorderLevel        int                          `json:"reorder_level"`
	ReorderQuantity     int                          `json:"reorder_quantity"`
	MinStockLevel       int                          `json:"min_stock_level"`
	MaxStockLevel       int                          `json:"max_stock_level"`
}

func (s *ProductService) CreateProduct(ctx context.Context, orgID primitive.ObjectID, req CreateProductRequest, createdBy primitive.ObjectID) (*models.Product, error) {
	// Check if SKU exists
	exists, err := s.productRepo.SKUExists(ctx, orgID, req.SKU, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check SKU existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("SKU '%s' already exists", req.SKU)
	}

	product := &models.Product{
		OrganizationID:     orgID,
		SKU:                req.SKU,
		Barcode:            req.Barcode,
		Name:               req.Name,
		Description:        req.Description,
		Type:               req.Type,
		TrackBatches:       req.TrackBatches,
		TrackSerialNumbers: req.TrackSerialNumbers,
		ValuationMethod:    req.ValuationMethod,
		BaseUOM:            req.BaseUOM,
		Weight:             req.Weight,
		CostPrice:          req.CostPrice,
		SellingPrice:       req.SellingPrice,
		ReorderLevel:       req.ReorderLevel,
		ReorderQuantity:    req.ReorderQuantity,
		MinStockLevel:      req.MinStockLevel,
		MaxStockLevel:      req.MaxStockLevel,
		Currency:           "USD",
	}

	if req.CategoryID != "" {
		catID, _ := primitive.ObjectIDFromHex(req.CategoryID)
		product.CategoryID = catID
	}

	product.BaseModel.CreatedBy = createdBy

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Product, error) {
	products, err := s.productRepo.FindAll(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id primitive.ObjectID, req CreateProductRequest, updatedBy primitive.ObjectID) (*models.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Check SKU uniqueness if changed
	if req.SKU != product.SKU {
		exists, err := s.productRepo.SKUExists(ctx, product.OrganizationID, req.SKU, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check SKU existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("SKU '%s' already exists", req.SKU)
		}
	}

	// Update fields
	product.SKU = req.SKU
	product.Barcode = req.Barcode
	product.Name = req.Name
	product.Description = req.Description
	product.Type = req.Type
	product.CostPrice = req.CostPrice
	product.SellingPrice = req.SellingPrice
	product.ReorderLevel = req.ReorderLevel
	product.ReorderQuantity = req.ReorderQuantity
	product.MinStockLevel = req.MinStockLevel
	product.MaxStockLevel = req.MaxStockLevel
	product.BaseModel.UpdatedBy = updatedBy

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	if err := s.productRepo.SoftDelete(ctx, id, deletedBy); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
