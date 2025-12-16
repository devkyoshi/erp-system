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
func (s *CategoryService) CreateCategory(ctx context.Context, category *models.ProductCategory, userOrgID primitive.ObjectID) error {
	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, category.OrganizationID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization not found")
	}

	// Check if user belongs to this organization
	if category.OrganizationID != userOrgID {
		return fmt.Errorf("unauthorized: cannot create category for different organization")
	}

	// Validate parent category if provided
	if category.ParentID != nil {
		parent, err := s.categoryRepo.FindByID(ctx, *category.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return fmt.Errorf("parent category not found")
		}
		if parent.OrganizationID != category.OrganizationID {
			return fmt.Errorf("parent category must belong to the same organization")
		}
	}

	// Check if name already exists at the same level
	exists, err = s.categoryRepo.CheckNameExists(ctx, category.OrganizationID, category.Name, category.ParentID, nil)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("category name already exists at this level")
	}

	return s.categoryRepo.Create(ctx, category)
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
func (s *CategoryService) ListCategories(ctx context.Context, orgID primitive.ObjectID, parentID *primitive.ObjectID, level *int, isActive *bool, page, limit int) ([]*models.ProductCategory, int64, error) {
	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, fmt.Errorf("organization not found")
	}

	return s.categoryRepo.FindByOrganization(ctx, orgID, parentID, level, isActive, page, limit)
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
func (s *CategoryService) UpdateCategory(ctx context.Context, id primitive.ObjectID, updates *models.ProductCategory, userOrgID primitive.ObjectID) error {
	// Get existing category
	existing, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("category not found")
	}

	// Check if user belongs to this organization
	if existing.OrganizationID != userOrgID {
		return fmt.Errorf("unauthorized: category belongs to different organization")
	}

	// Prevent changing organization
	if updates.OrganizationID != existing.OrganizationID {
		return fmt.Errorf("cannot change category organization")
	}

	// Validate parent category if changed
	if updates.ParentID != nil {
		// Cannot set self as parent
		if *updates.ParentID == id {
			return fmt.Errorf("category cannot be its own parent")
		}

		parent, err := s.categoryRepo.FindByID(ctx, *updates.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return fmt.Errorf("parent category not found")
		}
		if parent.OrganizationID != updates.OrganizationID {
			return fmt.Errorf("parent category must belong to the same organization")
		}

		// Check for circular reference (parent cannot be a descendant)
		if s.isDescendant(ctx, id, *updates.ParentID) {
			return fmt.Errorf("circular reference detected: parent cannot be a descendant")
		}
	}

	// Check if name already exists at the same level (excluding current category)
	if updates.Name != existing.Name || updates.ParentID != existing.ParentID {
		exists, err := s.categoryRepo.CheckNameExists(ctx, updates.OrganizationID, updates.Name, updates.ParentID, &id)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("category name already exists at this level")
		}
	}

	// Preserve certain fields
	updates.ID = existing.ID
	updates.CreatedAt = existing.CreatedAt
	updates.ProductCount = existing.ProductCount

	return s.categoryRepo.Update(ctx, updates)
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
