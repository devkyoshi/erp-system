package service

import (
	"context"
	"fmt"

	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	orgRepo      *repository.OrganizationRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, orgRepo *repository.OrganizationRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		orgRepo:      orgRepo,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, req CreateCategoryRequest, userOrgID primitive.ObjectID) (*models.ProductCategory, error) {
	var orgID primitive.ObjectID

	// Parse parent ID if provided
	var parentID *primitive.ObjectID
	if req.ParentID != nil && *req.ParentID != "" {
		pid, err := primitive.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent ID: %w", err)
		}
		parentID = &pid
	}

	// If parent exists, inherit organization from parent and verify access
	if parentID != nil {
		parent, err := s.categoryRepo.FindByID(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("parent category not found")
		}

		// Check if user belongs to parent's organization
		if parent.OrganizationID != userOrgID {
			return nil, fmt.Errorf("unauthorized: cannot create subcategory for parent in different organization")
		}

		// Inherit organization from parent
		orgID = parent.OrganizationID
	} else {
		// For root categories, use organization ID from request
		if req.OrganizationID == "" {
			return nil, fmt.Errorf("organization_id is required for root categories")
		}

		var err error
		orgID, err = primitive.ObjectIDFromHex(req.OrganizationID)
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
			return nil, fmt.Errorf("unauthorized: cannot create category for different organization")
		}
	}

	// Convert DTO to model
	category := &models.ProductCategory{
		OrganizationID: orgID,
		ParentID:       parentID,
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		IsActive:       req.IsActive,
		Metadata:       req.Metadata,
	}

	// Check if name already exists at the same level
	nameExists, err := s.categoryRepo.CheckNameExists(ctx, category.OrganizationID, category.Name, category.ParentID, nil)
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, fmt.Errorf("category name already exists at this level")
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(ctx context.Context, id primitive.ObjectID, userOrgID primitive.ObjectID) (*models.ProductCategory, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	// Check if user belongs to this organization
	if category.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: category belongs to different organization")
	}

	return category, nil
}

// ListCategories retrieves categories with filters
func (s *CategoryService) ListCategories(ctx context.Context, orgID primitive.ObjectID, parentID *primitive.ObjectID, level *int, isActive *bool, page, limit int) ([]*CategoryWithSubcategories, int64, error) {
	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, fmt.Errorf("organization not found")
	}

	categories, total, err := s.categoryRepo.FindByOrganization(ctx, orgID, parentID, level, isActive, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Fetch subcategories for each category
	result := make([]*CategoryWithSubcategories, len(categories))
	for i, cat := range categories {
		subcategories, err := s.categoryRepo.FindChildren(ctx, cat.ID)
		if err != nil {
			return nil, 0, err
		}

		result[i] = &CategoryWithSubcategories{
			ProductCategory: cat,
			Subcategories:   subcategories,
		}
	}

	return result, total, nil
}

// GetCategoryTree retrieves the full category tree
func (s *CategoryService) GetCategoryTree(ctx context.Context, orgID primitive.ObjectID) (interface{}, error) {
	categories, err := s.categoryRepo.GetCategoryTree(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Build hierarchical tree
	tree := s.buildTree(categories, nil)
	return tree, nil
}

// buildTree recursively builds a category tree
func (s *CategoryService) buildTree(categories []*models.ProductCategory, parentID *primitive.ObjectID) []map[string]interface{} {
	var result []map[string]interface{}

	for _, cat := range categories {
		if (parentID == nil && cat.ParentID == nil) || (parentID != nil && cat.ParentID != nil && *cat.ParentID == *parentID) {
			node := map[string]interface{}{
				"id":            cat.ID,
				"name":          cat.Name,
				"code":          cat.Code,
				"description":   cat.Description,
				"level":         cat.Level,
				"path":          cat.Path,
				"is_active":     cat.IsActive,
				"product_count": cat.ProductCount,
				"created_at":    cat.CreatedAt,
				"updated_at":    cat.UpdatedAt,
				"children":      s.buildTree(categories, &cat.ID),
			}
			result = append(result, node)
		}
	}

	return result
}

// GetRootCategories retrieves all root-level categories
func (s *CategoryService) GetRootCategories(ctx context.Context, orgID primitive.ObjectID) ([]*models.ProductCategory, error) {
	return s.categoryRepo.FindRootCategories(ctx, orgID)
}

// GetChildren retrieves direct children of a category
func (s *CategoryService) GetChildren(ctx context.Context, parentID primitive.ObjectID, userOrgID primitive.ObjectID) ([]*models.ProductCategory, error) {
	// Verify parent exists and user has access
	parent, err := s.categoryRepo.FindByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, fmt.Errorf("parent category not found")
	}
	if parent.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: category belongs to different organization")
	}

	return s.categoryRepo.FindChildren(ctx, parentID)
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(ctx context.Context, id primitive.ObjectID, req UpdateCategoryRequest, userOrgID primitive.ObjectID) (*models.ProductCategory, error) {
	// Get existing category
	existing, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("category not found")
	}

	// Check if user belongs to this organization
	if existing.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: category belongs to different organization")
	}

	// Apply updates
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Code != nil {
		existing.Code = *req.Code
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.Metadata != nil {
		existing.Metadata = req.Metadata
	}

	// Handle parent ID update
	if req.ParentID != nil {
		if *req.ParentID == "" {
			existing.ParentID = nil
		} else {
			parentID, err := primitive.ObjectIDFromHex(*req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("invalid parent ID: %w", err)
			}

			// Cannot set self as parent
			if parentID == id {
				return nil, fmt.Errorf("category cannot be its own parent")
			}

			parent, err := s.categoryRepo.FindByID(ctx, parentID)
			if err != nil {
				return nil, err
			}
			if parent == nil {
				return nil, fmt.Errorf("parent category not found")
			}
			if parent.OrganizationID != existing.OrganizationID {
				return nil, fmt.Errorf("parent category must belong to the same organization")
			}

			// Check for circular reference (parent cannot be a descendant)
			if s.isDescendant(ctx, id, parentID) {
				return nil, fmt.Errorf("circular reference detected: parent cannot be a descendant")
			}

			existing.ParentID = &parentID
		}
	}

	// Check if name already exists at the same level (excluding current category)
	if req.Name != nil || req.ParentID != nil {
		exists, err := s.categoryRepo.CheckNameExists(ctx, existing.OrganizationID, existing.Name, existing.ParentID, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("category name already exists at this level")
		}
	}

	if err := s.categoryRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// isDescendant checks if potentialDescendant is a descendant of categoryID
func (s *CategoryService) isDescendant(ctx context.Context, categoryID, potentialDescendant primitive.ObjectID) bool {
	current := potentialDescendant

	for {
		cat, err := s.categoryRepo.FindByID(ctx, current)
		if err != nil || cat == nil || cat.ParentID == nil {
			return false
		}

		if *cat.ParentID == categoryID {
			return true
		}

		current = *cat.ParentID
	}
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id primitive.ObjectID, userOrgID primitive.ObjectID) error {
	// Get category
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return fmt.Errorf("category not found")
	}

	// Check if user belongs to this organization
	if category.OrganizationID != userOrgID {
		return fmt.Errorf("unauthorized: category belongs to different organization")
	}

	// Check if category has children
	hasChildren, err := s.categoryRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return fmt.Errorf("cannot delete category with subcategories")
	}

	// Check if category has products
	hasProducts, err := s.categoryRepo.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return fmt.Errorf("cannot delete category with products")
	}

	return s.categoryRepo.Delete(ctx, id)
}

// UpdateProductCount updates the product count for a category
func (s *CategoryService) UpdateProductCount(ctx context.Context, categoryID primitive.ObjectID) error {
	return s.categoryRepo.UpdateProductCount(ctx, categoryID)
}

// Request DTOs
type CreateCategoryRequest struct {
	OrganizationID string                 `json:"organization_id"` // Required only for root categories, inherited from parent for subcategories
	Name           string                 `json:"name" binding:"required"`
	Code           string                 `json:"code"`
	Description    string                 `json:"description"`
	ParentID       *string                `json:"parent_id"` // If provided, organization_id is inherited from parent
	IsActive       bool                   `json:"is_active"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type UpdateCategoryRequest struct {
	Name        *string                `json:"name"`
	Code        *string                `json:"code"`
	Description *string                `json:"description"`
	ParentID    *string                `json:"parent_id"`
	IsActive    *bool                  `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CategoryFilter struct {
	ParentID *primitive.ObjectID
	Level    *int
	IsActive *bool
	Page     int
	Limit    int
}

type CategoryTreeNode struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Code         string             `json:"code"`
	Description  string             `json:"description"`
	Level        int                `json:"level"`
	Path         string             `json:"path"`
	IsActive     bool               `json:"is_active"`
	ProductCount int                `json:"product_count"`
	Children     []CategoryTreeNode `json:"children"`
}

// CategoryWithSubcategories represents a category with its subcategories
type CategoryWithSubcategories struct {
	*models.ProductCategory
	Subcategories []*models.ProductCategory `json:"subcategories"`
}
