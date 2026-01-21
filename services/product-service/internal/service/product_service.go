package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	brandRepo    *repository.BrandRepository
	orgRepo      *repository.OrganizationRepository
	unitRepo     *repository.UnitRepository
}

func NewProductService(
	productRepo *repository.ProductRepository,
	categoryRepo *repository.CategoryRepository,
	brandRepo *repository.BrandRepository,
	orgRepo *repository.OrganizationRepository,
	unitRepo *repository.UnitRepository,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		orgRepo:      orgRepo,
		unitRepo:     unitRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req CreateProductRequest, userOrgID primitive.ObjectID) (*models.Product, error) {
	// Parse organization ID
	orgID, err := primitive.ObjectIDFromHex(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization not found")
	}

	// Check if user belongs to this organization
	if orgID != userOrgID {
		return nil, fmt.Errorf("unauthorized: cannot create product for different organization")
	}

	// Convert DTO to model
	product, err := s.createProductRequestToModel(req, orgID)
	if err != nil {
		return nil, err
	}

	// Check if SKU already exists
	exists, err = s.productRepo.CheckSKUExists(ctx, product.OrganizationID, product.SKU, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("SKU already exists")
	}

	// Validate category if provided
	if !product.CategoryID.IsZero() {
		category, err := s.categoryRepo.FindByID(ctx, product.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			// Try to find if this is a subcategory ID
			parentCategory, err := s.categoryRepo.FindParentBySubcategoryID(ctx, product.CategoryID)
			if err != nil {
				return nil, err
			}
			if parentCategory != nil {
				// User passed a subcategory ID as category ID
				// Set CategoryID to parent, and SubcategoryID to the input ID
				product.SubcategoryID = product.CategoryID
				product.CategoryID = parentCategory.ID
				category = parentCategory
			} else {
				return nil, fmt.Errorf("category not found")
			}
		}
		if category.OrganizationID != product.OrganizationID {
			return nil, fmt.Errorf("category must belong to the same organization")
		}

		// Validate subcategory if provided
		if !product.SubcategoryID.IsZero() {
			subcategoryFound := false
			for _, sub := range category.Subcategories {
				if sub.ID == product.SubcategoryID && sub.DeletedAt == nil {
					subcategoryFound = true
					break
				}
			}
			if !subcategoryFound {
				return nil, fmt.Errorf("subcategory not found in the specified category")
			}
		}
	} else if !product.SubcategoryID.IsZero() {
		return nil, fmt.Errorf("subcategory_id requires category_id to be specified")
	}

	// Validate brand if provided
	if !product.BrandID.IsZero() {
		brand, err := s.brandRepo.FindByID(ctx, product.BrandID)
		if err != nil {
			return nil, err
		}
		if brand == nil {
			return nil, fmt.Errorf("brand not found")
		}
		if brand.OrganizationID != product.OrganizationID {
			return nil, fmt.Errorf("brand must belong to the same organization")
		}
	}

	// Set default values
	if product.Status == "" {
		product.Status = models.ProductStatusActive
	}
	if product.Type == "" {
		product.Type = models.ProductTypeFinished
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

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID, including its category details
func (s *ProductService) GetProduct(ctx context.Context, id primitive.ObjectID, userOrgID primitive.ObjectID) (*ProductDetailResponse, error) {
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

	response := &ProductDetailResponse{
		Product: *product,
	}

	// If category ID exists, fetch the category
	if !product.CategoryID.IsZero() {
		category, err := s.categoryRepo.FindByID(ctx, product.CategoryID)
		if err != nil {
			// Decide if you want to return the product anyway or fail
			// For now, we'll just log the error and return the product without category
			// In a real scenario, you might want to handle this more gracefully
			fmt.Printf("Warning: could not fetch category %s for product %s: %v\n", product.CategoryID.Hex(), product.ID.Hex(), err)
		}
		response.Category = category
	}

	// Populate unit details in location prices
	if err := s.populateUnits(ctx, &response.Product); err != nil {
		fmt.Printf("Warning: could not populate units for product %s: %v\n", product.ID.Hex(), err)
	}

	return response, nil
}

// ListProducts retrieves products with filters
func (s *ProductService) ListProducts(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*ProductListItemResponse, int64, error) {
	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, fmt.Errorf("organization not found")
	}

	products, total, err := s.productRepo.FindByOrganization(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Populate units for all products
	for _, product := range products {
		if err := s.populateUnits(ctx, product); err != nil {
			fmt.Printf("Warning: could not populate units for product %s: %v\n", product.ID.Hex(), err)
		}
	}

	// Populate categories and subcategories
	responses, err := s.populateCategories(ctx, products)
	if err != nil {
		return nil, 0, err
	}

	return responses, total, nil
}


// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id primitive.ObjectID, req UpdateProductRequest, userOrgID primitive.ObjectID) (*models.Product, error) {
	// Get existing product
	existing, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Check if user belongs to this organization
	if existing.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: product belongs to different organization")
	}

	// Apply updates to existing product
	if err := s.applyProductUpdates(ctx, existing, req); err != nil {
		return nil, err
	}

	// Check if SKU is being changed and if new SKU exists
	if req.SKU != nil && *req.SKU != existing.SKU {
		exists, err := s.productRepo.CheckSKUExists(ctx, existing.OrganizationID, *req.SKU, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("SKU already exists")
		}
	}

	// Validate category if provided and changed
	if req.CategoryID != nil {
		catID, err := primitive.ObjectIDFromHex(*req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		if catID != existing.CategoryID {
			category, err := s.categoryRepo.FindByID(ctx, catID)
			if err != nil {
				return nil, err
			}
			if category == nil {
				// Try to find if this is a subcategory ID
				parentCategory, err := s.categoryRepo.FindParentBySubcategoryID(ctx, catID)
				if err != nil {
					return nil, err
				}
				if parentCategory != nil {
					// User passed a subcategory ID as category ID
					existing.SubcategoryID = catID
					catID = parentCategory.ID
					category = parentCategory
				} else {
					return nil, fmt.Errorf("category not found")
				}
			}
			if category.OrganizationID != existing.OrganizationID {
				return nil, fmt.Errorf("category must belong to the same organization")
			}
			existing.CategoryID = catID
		}
	}


	// Validate subcategory if provided
	if req.SubcategoryID != nil {
		subID, err := primitive.ObjectIDFromHex(*req.SubcategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid subcategory ID: %w", err)
		}

		// Determine the final category ID (either updated or existing)
		finalCategoryID := existing.CategoryID
		if req.CategoryID != nil {
			finalCategoryID, _ = primitive.ObjectIDFromHex(*req.CategoryID)
		}

		if finalCategoryID.IsZero() {
			return nil, fmt.Errorf("subcategory_id requires category_id to be specified")
		}

		category, err := s.categoryRepo.FindByID(ctx, finalCategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, fmt.Errorf("category not found")
		}

		subcategoryFound := false
		for _, sub := range category.Subcategories {
			if sub.ID == subID && sub.DeletedAt == nil {
				subcategoryFound = true
				break
			}
		}
		if !subcategoryFound {
			return nil, fmt.Errorf("subcategory not found in the specified category")
		}
	}

	// Validate brand if provided and changed
	if req.BrandID != nil {
		brandID, err := primitive.ObjectIDFromHex(*req.BrandID)
		if err != nil {
			return nil, fmt.Errorf("invalid brand ID: %w", err)
		}
		if brandID != existing.BrandID {
			brand, err := s.brandRepo.FindByID(ctx, brandID)
			if err != nil {
				return nil, err
			}
			if brand == nil {
				return nil, fmt.Errorf("brand not found")
			}
			if brand.OrganizationID != existing.OrganizationID {
				return nil, fmt.Errorf("brand must belong to the same organization")
			}
		}
	}

	if err := s.productRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
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
	products, err := s.productRepo.GetLowStockProducts(ctx, orgID)
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		if err := s.populateUnits(ctx, product); err != nil {
			fmt.Printf("Warning: could not populate units for product %s: %v\n", product.ID.Hex(), err)
		}
	}

	return products, nil
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

	if err := s.populateUnits(ctx, product); err != nil {
		fmt.Printf("Warning: could not populate units for product %s: %v\n", product.ID.Hex(), err)
	}

	return product, nil
}

// Request DTOs
type CreateProductRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	SKU            string  `json:"sku" binding:"required"`
	Barcode        string  `json:"barcode"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Type           string  `json:"type"`
	Status         string  `json:"status"`
	CategoryID     *string `json:"category_id"`
	SubcategoryID  *string `json:"subcategory_id"` // NEW: Subcategory reference
	BrandID        *string `json:"brand_id"`
	ManufacturerID *string `json:"manufacturer_id"`

	// Inventory Settings
	TrackInventory     bool   `json:"track_inventory"`
	TrackBatches       bool   `json:"track_batches"`
	TrackSerialNumbers bool   `json:"track_serial_numbers"`
	ValuationMethod    string `json:"valuation_method"`

	// Unit of Measure
	BaseUnitID     *string  `json:"base_unit_id"`
	AllowedUnitIDs []string `json:"allowed_unit_ids"`

	// Dimensions & Weight
	Weight        float64 `json:"weight"`
	WeightUnit    string  `json:"weight_unit"`
	Length        float64 `json:"length"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
	DimensionUnit string  `json:"dimension_unit"`
	Volume        float64 `json:"volume"`
	VolumeUnit    string  `json:"volume_unit"`

	// Pricing - Location-wise
	LocationPrices []LocationPriceRequest `json:"location_prices"`

	// Tax & Accounting
	TaxCategoryID *string `json:"tax_category_id"`
	HSNCode       string  `json:"hsn_code"`
	SACCode       string  `json:"sac_code"`

	// Reorder Settings
	ReorderLevel    int `json:"reorder_level"`
	ReorderQuantity int `json:"reorder_quantity"`
	MinStockLevel   int `json:"min_stock_level"`
	MaxStockLevel   int `json:"max_stock_level"`
	SafetyStock     int `json:"safety_stock"`

	// Supplier Info
	DefaultSupplierID *string  `json:"default_supplier_id"`
	SupplierIDs       []string `json:"supplier_ids"`
	LeadTimeDays      int      `json:"lead_time_days"`

	// Quality & Expiry
	ShelfLifeDays int  `json:"shelf_life_days"`
	RequiresQC    bool `json:"requires_qc"`
	Perishable    bool `json:"perishable"`
	Hazardous     bool `json:"hazardous"`

	// Images & Attachments
	Images         []string          `json:"images"`
	Thumbnail      string            `json:"thumbnail"`
	Specifications map[string]string `json:"specifications"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata"`
}

type UpdateProductRequest struct {
	SKU            *string `json:"sku"`
	Barcode        *string `json:"barcode"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Type           *string `json:"type"`
	Status         *string `json:"status"`
	CategoryID     *string `json:"category_id"`
	SubcategoryID  *string `json:"subcategory_id"` // NEW: Subcategory reference
	BrandID        *string `json:"brand_id"`
	ManufacturerID *string `json:"manufacturer_id"`

	// Inventory Settings
	TrackInventory     *bool   `json:"track_inventory"`
	TrackBatches       *bool   `json:"track_batches"`
	TrackSerialNumbers *bool   `json:"track_serial_numbers"`
	ValuationMethod    *string `json:"valuation_method"`

	// Unit of Measure
	BaseUnitID     *string  `json:"base_unit_id"`
	AllowedUnitIDs []string `json:"allowed_unit_ids"`

	// Dimensions & Weight
	Weight        *float64 `json:"weight"`
	WeightUnit    *string  `json:"weight_unit"`
	Length        *float64 `json:"length"`
	Width         *float64 `json:"width"`
	Height        *float64 `json:"height"`
	DimensionUnit *string  `json:"dimension_unit"`
	Volume        *float64 `json:"volume"`
	VolumeUnit    *string  `json:"volume_unit"`

	// Pricing - Location-wise
	LocationPrices []LocationPriceRequest `json:"location_prices"`

	// Tax & Accounting
	TaxCategoryID *string `json:"tax_category_id"`
	HSNCode       *string `json:"hsn_code"`
	SACCode       *string `json:"sac_code"`

	// Reorder Settings
	ReorderLevel    *int `json:"reorder_level"`
	ReorderQuantity *int `json:"reorder_quantity"`
	MinStockLevel   *int `json:"min_stock_level"`
	MaxStockLevel   *int `json:"max_stock_level"`
	SafetyStock     *int `json:"safety_stock"`

	// Supplier Info
	DefaultSupplierID *string  `json:"default_supplier_id"`
	SupplierIDs       []string `json:"supplier_ids"`
	LeadTimeDays      *int     `json:"lead_time_days"`

	// Quality & Expiry
	ShelfLifeDays *int  `json:"shelf_life_days"`
	RequiresQC    *bool `json:"requires_qc"`
	Perishable    *bool `json:"perishable"`
	Hazardous     *bool `json:"hazardous"`

	// Images & Attachments
	Images         []string          `json:"images"`
	Thumbnail      *string           `json:"thumbnail"`
	Specifications map[string]string `json:"specifications"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata"`
}

// LocationPriceRequest represents location-wise pricing request
type LocationPriceRequest struct {
	LocationID     string  `json:"location_id" binding:"required"`
	LocationName   string  `json:"location_name"`
	PurchaseUnitID string  `json:"purchase_unit_id"`
	SellingUnitID  string  `json:"selling_unit_id"`
	CostPrice      float64 `json:"cost_price"`
	SellingPrice   float64 `json:"selling_price"`
	MRP            float64 `json:"mrp"`
	InitialStock   float64 `json:"initial_stock"`
	Currency       string  `json:"currency"`
	IsActive       bool    `json:"is_active"`
}

type ProductFilter struct {
	CategoryID     *primitive.ObjectID
	BrandID        *primitive.ObjectID
	Status         *models.ProductStatus
	Type           *models.ProductType
	TrackInventory *bool
	Search         string
	Page           int
	Limit          int
}


type ProductDetailResponse struct {
	models.Product
	Category *models.ProductCategory `json:"category,omitempty"`
}

// CategoryInfo represents simplified category information for list responses
type CategoryInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Code string             `json:"code"`
}

// SubcategoryInfo represents simplified subcategory information for list responses
type SubcategoryInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Code string             `json:"code"`
}

// ProductListItemResponse extends Product with category and subcategory details
type ProductListItemResponse struct {
	models.Product
	Category    *CategoryInfo    `json:"category,omitempty"`
	Subcategory *SubcategoryInfo `json:"subcategory,omitempty"`
}


type ProductResponse struct {
	ID             primitive.ObjectID   `json:"id"`
	OrganizationID primitive.ObjectID   `json:"organization_id"`
	SKU            string               `json:"sku"`
	Barcode        string               `json:"barcode"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Type           models.ProductType   `json:"type"`
	Status         models.ProductStatus `json:"status"`

	// Stock Info
	TotalStock     float64 `json:"total_stock"`
	AvailableStock float64 `json:"available_stock"`
	AllocatedStock float64 `json:"allocated_stock"`
	InTransitStock float64 `json:"in_transit_stock"`
	StockValue     float64 `json:"stock_value"`

	// Pricing - Location-wise
	LocationPrices []models.LocationPrice `json:"location_prices"`

	// References
	CategoryID *primitive.ObjectID `json:"category_id,omitempty"`
	BrandID    *primitive.ObjectID `json:"brand_id,omitempty"`

	// Timestamps
	CreatedAt primitive.DateTime `json:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
}

// Helper methods to convert between DTOs and models
func (s *ProductService) createProductRequestToModel(req CreateProductRequest, orgID primitive.ObjectID) (*models.Product, error) {
	product := &models.Product{
		OrganizationID:     orgID,
		SKU:                req.SKU,
		Barcode:            req.Barcode,
		Name:               req.Name,
		Description:        req.Description,
		Type:               models.ProductType(req.Type),
		Status:             models.ProductStatus(req.Status),
		TrackInventory:     req.TrackInventory,
		TrackBatches:       req.TrackBatches,
		TrackSerialNumbers: req.TrackSerialNumbers,
		ValuationMethod:    models.StockValuationMethod(req.ValuationMethod),
		Weight:             req.Weight,
		WeightUnit:         req.WeightUnit,
		Length:             req.Length,
		Width:              req.Width,
		Height:             req.Height,
		DimensionUnit:      req.DimensionUnit,
		Volume:             req.Volume,
		VolumeUnit:         req.VolumeUnit,
		HSNCode:            req.HSNCode,
		SACCode:            req.SACCode,
		ReorderLevel:       req.ReorderLevel,
		ReorderQuantity:    req.ReorderQuantity,
		MinStockLevel:      req.MinStockLevel,
		MaxStockLevel:      req.MaxStockLevel,
		SafetyStock:        req.SafetyStock,
		LeadTimeDays:       req.LeadTimeDays,
		ShelfLifeDays:      req.ShelfLifeDays,
		RequiresQC:         req.RequiresQC,
		Perishable:         req.Perishable,
		Hazardous:          req.Hazardous,
		Images:             req.Images,
		Thumbnail:          req.Thumbnail,
		Specifications:     req.Specifications,
		Metadata:           req.Metadata,
	}

	// Parse optional ObjectIDs
	if req.CategoryID != nil {
		id, err := primitive.ObjectIDFromHex(*req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		product.CategoryID = id
	}

	if req.SubcategoryID != nil {
		id, err := primitive.ObjectIDFromHex(*req.SubcategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid subcategory ID: %w", err)
		}
		product.SubcategoryID = id
	}

	if req.BrandID != nil {
		id, err := primitive.ObjectIDFromHex(*req.BrandID)
		if err != nil {
			return nil, fmt.Errorf("invalid brand ID: %w", err)
		}
		product.BrandID = id
	}

	//TODO: 
	// if req.ManufacturerID != nil {
	// 	id, err := primitive.ObjectIDFromHex(*req.ManufacturerID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid manufacturer ID: %w", err)
	// 	}
	// 	product.ManufacturerID = id
	// }

	if req.BaseUnitID != nil {
		id, err := primitive.ObjectIDFromHex(*req.BaseUnitID)
		if err != nil {
			return nil, fmt.Errorf("invalid base unit ID: %w", err)
		}
		product.BaseUnitID = id
	}

	// TODO: if req.TaxCategoryID != nil {
	// 	id, err := primitive.ObjectIDFromHex(*req.TaxCategoryID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid tax category ID: %w", err)
	// 	}
	// 	product.TaxCategoryID = id
	// }

	//TODO:
	// if req.DefaultSupplierID != nil {
	// 	id, err := primitive.ObjectIDFromHex(*req.DefaultSupplierID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid default supplier ID: %w", err)
	// 	}
	// 	product.DefaultSupplierID = id
	// }

	// Parse allowed unit IDs
	if len(req.AllowedUnitIDs) > 0 {
		allowedIDs := make([]primitive.ObjectID, 0, len(req.AllowedUnitIDs))
		for _, idStr := range req.AllowedUnitIDs {
			id, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				return nil, fmt.Errorf("invalid allowed unit ID %s: %w", idStr, err)
			}
			allowedIDs = append(allowedIDs, id)
		}
		product.AllowedUnitIDs = allowedIDs
	}

	// Parse supplier IDs
	if len(req.SupplierIDs) > 0 {
		supplierIDs := make([]primitive.ObjectID, 0, len(req.SupplierIDs))
		for _, idStr := range req.SupplierIDs {
			id, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				return nil, fmt.Errorf("invalid supplier ID %s: %w", idStr, err)
			}
			supplierIDs = append(supplierIDs, id)
		}
		product.SupplierIDs = supplierIDs
	}

	// Parse location prices
	if len(req.LocationPrices) > 0 {
		locationPrices := make([]models.LocationPrice, 0, len(req.LocationPrices))
		timestamp := primitive.DateTime(time.Now().Unix() * 1000)

		for _, lpReq := range req.LocationPrices {
			locID, err := primitive.ObjectIDFromHex(lpReq.LocationID)
			if err != nil {
				return nil, fmt.Errorf("invalid location ID %s: %w", lpReq.LocationID, err)
			}

			lp := models.LocationPrice{
				LocationID:   locID,
				LocationName: lpReq.LocationName,
				CostPrice:    lpReq.CostPrice,
				SellingPrice: lpReq.SellingPrice,
				MRP:          lpReq.MRP,
				InitialStock: lpReq.InitialStock,
				Currency:     lpReq.Currency,
				IsActive:     lpReq.IsActive,
				CreatedAt:    int64(timestamp),
				ModifiedAt:   int64(timestamp),
			}

			if lpReq.PurchaseUnitID != "" {
				unitID, err := primitive.ObjectIDFromHex(lpReq.PurchaseUnitID)
				if err != nil {
					return nil, fmt.Errorf("invalid purchase unit ID %s: %w", lpReq.PurchaseUnitID, err)
				}
				lp.PurchaseUnitID = unitID
			}

			if lpReq.SellingUnitID != "" {
				unitID, err := primitive.ObjectIDFromHex(lpReq.SellingUnitID)
				if err != nil {
					return nil, fmt.Errorf("invalid selling unit ID %s: %w", lpReq.SellingUnitID, err)
				}
				lp.SellingUnitID = unitID
			}

			// Set default currency if not provided
			if lp.Currency == "" {
				lp.Currency = "USD"
			}

			locationPrices = append(locationPrices, lp)
		}
		product.LocationPrices = locationPrices
	}

	return product, nil
}

func (s *ProductService) applyProductUpdates(ctx context.Context, product *models.Product, req UpdateProductRequest) error {
	if req.SKU != nil {
		product.SKU = *req.SKU
	}
	if req.Barcode != nil {
		product.Barcode = *req.Barcode
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Type != nil {
		product.Type = models.ProductType(*req.Type)
	}
	if req.Status != nil {
		product.Status = models.ProductStatus(*req.Status)
	}
	if req.CategoryID != nil {
		id, err := primitive.ObjectIDFromHex(*req.CategoryID)
		if err != nil {
			return fmt.Errorf("invalid category ID: %w", err)
		}
		product.CategoryID = id
	}
	if req.SubcategoryID != nil {
		id, err := primitive.ObjectIDFromHex(*req.SubcategoryID)
		if err != nil {
			return fmt.Errorf("invalid subcategory ID: %w", err)
		}
		product.SubcategoryID = id
	}
	if req.BrandID != nil {
		id, err := primitive.ObjectIDFromHex(*req.BrandID)
		if err != nil {
			return fmt.Errorf("invalid brand ID: %w", err)
		}
		product.BrandID = id
	}
	if req.ManufacturerID != nil {
		id, err := primitive.ObjectIDFromHex(*req.ManufacturerID)
		if err != nil {
			return fmt.Errorf("invalid manufacturer ID: %w", err)
		}
		product.ManufacturerID = id
	}
	if req.TrackInventory != nil {
		product.TrackInventory = *req.TrackInventory
	}
	if req.TrackBatches != nil {
		product.TrackBatches = *req.TrackBatches
	}
	if req.TrackSerialNumbers != nil {
		product.TrackSerialNumbers = *req.TrackSerialNumbers
	}
	if req.ValuationMethod != nil {
		product.ValuationMethod = models.StockValuationMethod(*req.ValuationMethod)
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.WeightUnit != nil {
		product.WeightUnit = *req.WeightUnit
	}
	if req.Length != nil {
		product.Length = *req.Length
	}
	if req.Width != nil {
		product.Width = *req.Width
	}
	if req.Height != nil {
		product.Height = *req.Height
	}
	if req.DimensionUnit != nil {
		product.DimensionUnit = *req.DimensionUnit
	}
	if req.Volume != nil {
		product.Volume = *req.Volume
	}
	if req.VolumeUnit != nil {
		product.VolumeUnit = *req.VolumeUnit
	}
	if req.HSNCode != nil {
		product.HSNCode = *req.HSNCode
	}
	if req.SACCode != nil {
		product.SACCode = *req.SACCode
	}
	if req.ReorderLevel != nil {
		product.ReorderLevel = *req.ReorderLevel
	}
	if req.ReorderQuantity != nil {
		product.ReorderQuantity = *req.ReorderQuantity
	}
	if req.MinStockLevel != nil {
		product.MinStockLevel = *req.MinStockLevel
	}
	if req.MaxStockLevel != nil {
		product.MaxStockLevel = *req.MaxStockLevel
	}
	if req.SafetyStock != nil {
		product.SafetyStock = *req.SafetyStock
	}
	if req.LeadTimeDays != nil {
		product.LeadTimeDays = *req.LeadTimeDays
	}
	if req.ShelfLifeDays != nil {
		product.ShelfLifeDays = *req.ShelfLifeDays
	}
	if req.RequiresQC != nil {
		product.RequiresQC = *req.RequiresQC
	}
	if req.Perishable != nil {
		product.Perishable = *req.Perishable
	}
	if req.Hazardous != nil {
		product.Hazardous = *req.Hazardous
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.Thumbnail != nil {
		product.Thumbnail = *req.Thumbnail
	}
	if req.Specifications != nil {
		product.Specifications = req.Specifications
	}
	if req.Metadata != nil {
		product.Metadata = req.Metadata
	}

	// Handle location prices
	if len(req.LocationPrices) > 0 {
		locationPrices := make([]models.LocationPrice, 0, len(req.LocationPrices))
		timestamp := primitive.DateTime(time.Now().Unix() * 1000)

		for _, lpReq := range req.LocationPrices {
			locID, err := primitive.ObjectIDFromHex(lpReq.LocationID)
			if err != nil {
				return fmt.Errorf("invalid location ID %s: %w", lpReq.LocationID, err)
			}

			lp := models.LocationPrice{
				LocationID:   locID,
				LocationName: lpReq.LocationName,
				CostPrice:    lpReq.CostPrice,
				SellingPrice: lpReq.SellingPrice,
				MRP:          lpReq.MRP,
				InitialStock: lpReq.InitialStock,
				Currency:     lpReq.Currency,
				IsActive:     lpReq.IsActive,
				CreatedAt:    int64(timestamp),
				ModifiedAt:   int64(timestamp),
			}

			if lpReq.PurchaseUnitID != "" {
				unitID, err := primitive.ObjectIDFromHex(lpReq.PurchaseUnitID)
				if err != nil {
					return fmt.Errorf("invalid purchase unit ID %s: %w", lpReq.PurchaseUnitID, err)
				}
				lp.PurchaseUnitID = unitID
			}

			if lpReq.SellingUnitID != "" {
				unitID, err := primitive.ObjectIDFromHex(lpReq.SellingUnitID)
				if err != nil {
					return fmt.Errorf("invalid selling unit ID %s: %w", lpReq.SellingUnitID, err)
				}
				lp.SellingUnitID = unitID
			}

			// Set default currency if not provided
			if lp.Currency == "" {
				lp.Currency = "USD"
			}

			locationPrices = append(locationPrices, lp)
		}
		product.LocationPrices = locationPrices
	}

	return nil
}

// Helper: populateUnits fetches and attaches Unit details to LocationPrices
func (s *ProductService) populateUnits(ctx context.Context, product *models.Product) error {
	if len(product.LocationPrices) == 0 {
		return nil
	}

	// Collect unique unit IDs
	unitIDs := make([]primitive.ObjectID, 0)
	unitIDMap := make(map[primitive.ObjectID]bool)

	for _, lp := range product.LocationPrices {
		if !lp.PurchaseUnitID.IsZero() {
			if !unitIDMap[lp.PurchaseUnitID] {
				unitIDs = append(unitIDs, lp.PurchaseUnitID)
				unitIDMap[lp.PurchaseUnitID] = true
			}
		}
		if !lp.SellingUnitID.IsZero() {
			if !unitIDMap[lp.SellingUnitID] {
				unitIDs = append(unitIDs, lp.SellingUnitID)
				unitIDMap[lp.SellingUnitID] = true
			}
		}
	}

	if len(unitIDs) == 0 {
		return nil
	}

	// Fetch units
	units, err := s.unitRepo.FindByIDs(ctx, unitIDs)
	if err != nil {
		return err
	}

	// Map units for quick lookup
	unitsMap := make(map[primitive.ObjectID]*models.Unit)
	for _, unit := range units {
		unitsMap[unit.ID] = unit
	}

	// Assign units to location prices
	for i := range product.LocationPrices {
		if !product.LocationPrices[i].PurchaseUnitID.IsZero() {
			if unit, ok := unitsMap[product.LocationPrices[i].PurchaseUnitID]; ok {
				product.LocationPrices[i].PurchaseUnit = unit
			}
		}
		if !product.LocationPrices[i].SellingUnitID.IsZero() {
			if unit, ok := unitsMap[product.LocationPrices[i].SellingUnitID]; ok {
				product.LocationPrices[i].SellingUnit = unit
			}
		}
	}

	return nil
}

// Helper: populateCategories fetches and attaches category and subcategory details to products
func (s *ProductService) populateCategories(ctx context.Context, products []*models.Product) ([]*ProductListItemResponse, error) {
	if len(products) == 0 {
		return []*ProductListItemResponse{}, nil
	}

	// Collect unique category IDs
	categoryIDs := make([]primitive.ObjectID, 0)
	categoryIDMap := make(map[primitive.ObjectID]bool)

	for _, product := range products {
		if !product.CategoryID.IsZero() && !categoryIDMap[product.CategoryID] {
			categoryIDs = append(categoryIDs, product.CategoryID)
			categoryIDMap[product.CategoryID] = true
		}
	}

	// Fetch all categories in one query
	categoriesMap := make(map[primitive.ObjectID]*models.ProductCategory)
	if len(categoryIDs) > 0 {
		for _, catID := range categoryIDs {
			category, err := s.categoryRepo.FindByID(ctx, catID)
			if err != nil {
				fmt.Printf("Warning: could not fetch category %s: %v\n", catID.Hex(), err)
				continue
			}
			if category != nil {
				categoriesMap[catID] = category
			}
		}
	}

	// Build response with category and subcategory details
	responses := make([]*ProductListItemResponse, 0, len(products))
	for _, product := range products {
		response := &ProductListItemResponse{
			Product: *product,
		}

		// Add category info if exists
		if !product.CategoryID.IsZero() {
			if category, ok := categoriesMap[product.CategoryID]; ok {
				response.Category = &CategoryInfo{
					ID:   category.ID,
					Name: category.Name,
					Code: category.Code,
				}

				// Add subcategory info if exists
				if !product.SubcategoryID.IsZero() {
					for _, subcat := range category.Subcategories {
						if subcat.ID == product.SubcategoryID && subcat.DeletedAt == nil {
							response.Subcategory = &SubcategoryInfo{
								ID:   subcat.ID,
								Name: subcat.Name,
								Code: subcat.Code,
							}
							break
						}
					}
				}
			}
		}

		responses = append(responses, response)
	}

	return responses, nil
}
