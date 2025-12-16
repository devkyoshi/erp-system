package service

import (
	"context"
	"fmt"

	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	brandRepo    *repository.BrandRepository
	orgRepo      *repository.OrganizationRepository
}

func NewProductService(
	productRepo *repository.ProductRepository,
	categoryRepo *repository.CategoryRepository,
	brandRepo *repository.BrandRepository,
	orgRepo *repository.OrganizationRepository,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		orgRepo:      orgRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product, userOrgID primitive.ObjectID) error {
	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, product.OrganizationID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization not found")
	}

	// Check if user belongs to this organization
	if product.OrganizationID != userOrgID {
		return fmt.Errorf("unauthorized: cannot create product for different organization")
	}

	// Check if SKU already exists
	exists, err = s.productRepo.CheckSKUExists(ctx, product.OrganizationID, product.SKU, nil)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("SKU already exists")
	}

	// Validate category if provided
	if !product.CategoryID.IsZero() {
		category, err := s.categoryRepo.FindByID(ctx, product.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return fmt.Errorf("category not found")
		}
		if category.OrganizationID != product.OrganizationID {
			return fmt.Errorf("category must belong to the same organization")
		}
	}

	// Validate brand if provided
	if !product.BrandID.IsZero() {
		brand, err := s.brandRepo.FindByID(ctx, product.BrandID)
		if err != nil {
			return err
		}
		if brand == nil {
			return fmt.Errorf("brand not found")
		}
		if brand.OrganizationID != product.OrganizationID {
			return fmt.Errorf("brand must belong to the same organization")
		}
	}

	// Set default values
	if product.Status == "" {
		product.Status = models.ProductStatusActive
	}
	if product.Type == "" {
		product.Type = models.ProductTypeFinished
	}
	if product.Currency == "" {
		product.Currency = "USD"
	}
	if product.ValuationMethod == "" {
		product.ValuationMethod = models.ValuationFIFO
	}

	// Initialize stock levels
	product.TotalStock = 0
	product.AvailableStock = 0
	product.AllocatedStock = 0
	product.InTransitStock = 0
	product.StockValue = 0
	product.TotalSold = 0
	product.TotalPurchased = 0

	return s.productRepo.Create(ctx, product)
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id primitive.ObjectID, userOrgID primitive.ObjectID) (*models.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Check if user belongs to this organization
	if product.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: product belongs to different organization")
	}

	return product, nil
}

// ListProducts retrieves products with filters
func (s *ProductService) ListProducts(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Product, int64, error) {
	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, fmt.Errorf("organization not found")
	}

	return s.productRepo.FindByOrganization(ctx, orgID, filters, page, limit)
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id primitive.ObjectID, updates *models.Product, userOrgID primitive.ObjectID) error {
	// Get existing product
	existing, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	// Check if user belongs to this organization
	if existing.OrganizationID != userOrgID {
		return fmt.Errorf("unauthorized: product belongs to different organization")
	}

	// Prevent changing organization
	if updates.OrganizationID != existing.OrganizationID {
		return fmt.Errorf("cannot change product organization")
	}

	// Check if SKU is being changed and if new SKU exists
	if updates.SKU != existing.SKU {
		exists, err := s.productRepo.CheckSKUExists(ctx, updates.OrganizationID, updates.SKU, &id)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("SKU already exists")
		}
	}

	// Validate category if provided and changed
	if !updates.CategoryID.IsZero() && updates.CategoryID != existing.CategoryID {
		category, err := s.categoryRepo.FindByID(ctx, updates.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return fmt.Errorf("category not found")
		}
		if category.OrganizationID != updates.OrganizationID {
			return fmt.Errorf("category must belong to the same organization")
		}
	}

	// Validate brand if provided and changed
	if !updates.BrandID.IsZero() && updates.BrandID != existing.BrandID {
		brand, err := s.brandRepo.FindByID(ctx, updates.BrandID)
		if err != nil {
			return err
		}
		if brand == nil {
			return fmt.Errorf("brand not found")
		}
		if brand.OrganizationID != updates.OrganizationID {
			return fmt.Errorf("brand must belong to the same organization")
		}
	}

	// Preserve certain fields
	updates.ID = existing.ID
	updates.CreatedAt = existing.CreatedAt
	updates.TotalStock = existing.TotalStock
	updates.AvailableStock = existing.AvailableStock
	updates.AllocatedStock = existing.AllocatedStock
	updates.InTransitStock = existing.InTransitStock
	updates.StockValue = existing.StockValue
	updates.TotalSold = existing.TotalSold
	updates.TotalPurchased = existing.TotalPurchased
	updates.LastSoldDate = existing.LastSoldDate
	updates.LastPurchaseDate = existing.LastPurchaseDate

	return s.productRepo.Update(ctx, updates)
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id primitive.ObjectID, userOrgID primitive.ObjectID) error {
	// Get product
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("product not found")
	}

	// Check if user belongs to this organization
	if product.OrganizationID != userOrgID {
		return fmt.Errorf("unauthorized: product belongs to different organization")
	}

	// Optional: Check if product has stock
	if product.TotalStock > 0 {
		return fmt.Errorf("cannot delete product with existing stock")
	}

	return s.productRepo.Delete(ctx, id)
}

// GetLowStockProducts retrieves products below reorder level
func (s *ProductService) GetLowStockProducts(ctx context.Context, orgID primitive.ObjectID) ([]*models.Product, error) {
	return s.productRepo.GetLowStockProducts(ctx, orgID)
}

// GetProductBySKU retrieves a product by SKU
func (s *ProductService) GetProductBySKU(ctx context.Context, orgID primitive.ObjectID, sku string, userOrgID primitive.ObjectID) (*models.Product, error) {
	// Check if user belongs to this organization
	if orgID != userOrgID {
		return nil, fmt.Errorf("unauthorized: cannot access products from different organization")
	}

	product, err := s.productRepo.FindBySKU(ctx, orgID, sku)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}