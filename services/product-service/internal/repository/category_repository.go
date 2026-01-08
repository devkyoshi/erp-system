package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository struct {
	collection        *mongo.Collection
	productCollection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{
		collection:        db.Collection("categories"),
		productCollection: db.Collection("products"),
	}
}

// Create inserts a new category
func (r *CategoryRepository) Create(ctx context.Context, category *models.ProductCategory) error {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.ProductCount = 0

	// Build path
	if category.ParentID != nil {
		parent, err := r.FindByID(ctx, *category.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return fmt.Errorf("parent category not found")
		}
		category.Level = parent.Level + 1
		category.Path = parent.Path + "/" + strings.ToLower(category.Name)
	} else {
		category.Level = 0
		category.Path = "/" + strings.ToLower(category.Name)
	}

	_, err := r.collection.InsertOne(ctx, category)
	return err
}

// FindByID retrieves a category by ID
func (r *CategoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.ProductCategory, error) {
	var category models.ProductCategory
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindByOrganization retrieves categories for an organization with optional filters
func (r *CategoryRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, parentID *primitive.ObjectID, level *int, isActive *bool, page, limit int) ([]*models.ProductCategory, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	if parentID != nil {
		filter["parent_id"] = *parentID
	}

	if level != nil {
		filter["level"] = *level
	}

	if isActive != nil {
		filter["is_active"] = *isActive
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Setup pagination
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "level", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// FindRootCategories retrieves all root-level categories (level 0)
func (r *CategoryRepository) FindRootCategories(ctx context.Context, orgID primitive.ObjectID) ([]*models.ProductCategory, error) {
	filter := bson.M{
		"organization_id": orgID,
		"level":           0,
		"deleted_at":      nil,
	}

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// FindChildren retrieves direct children of a category
func (r *CategoryRepository) FindChildren(ctx context.Context, parentID primitive.ObjectID) ([]*models.ProductCategory, error) {
	filter := bson.M{
		"parent_id":  parentID,
		"deleted_at": nil,
	}

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// FindByPath retrieves a category by its path
func (r *CategoryRepository) FindByPath(ctx context.Context, orgID primitive.ObjectID, path string) (*models.ProductCategory, error) {
	var category models.ProductCategory
	filter := bson.M{
		"organization_id": orgID,
		"path":            path,
		"deleted_at":      nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(ctx context.Context, category *models.ProductCategory) error {
	category.UpdatedAt = time.Now()

	// If name changed, update path for this category and all descendants
	oldCategory, err := r.FindByID(ctx, category.ID)
	if err != nil {
		return err
	}
	if oldCategory == nil {
		return fmt.Errorf("category not found")
	}

	// Rebuild path if name or parent changed
	if oldCategory.Name != category.Name || (oldCategory.ParentID != category.ParentID) {
		if category.ParentID != nil {
			parent, err := r.FindByID(ctx, *category.ParentID)
			if err != nil {
				return err
			}
			if parent == nil {
				return fmt.Errorf("parent category not found")
			}
			category.Level = parent.Level + 1
			category.Path = parent.Path + "/" + strings.ToLower(category.Name)
		} else {
			category.Level = 0
			category.Path = "/" + strings.ToLower(category.Name)
		}

		// Update paths of all descendants
		if err := r.updateDescendantPaths(ctx, category.ID, oldCategory.Path, category.Path); err != nil {
			return err
		}
	}

	filter := bson.M{
		"_id":        category.ID,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": category,
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// updateDescendantPaths updates paths of all descendant categories
func (r *CategoryRepository) updateDescendantPaths(ctx context.Context, categoryID primitive.ObjectID, oldPath, newPath string) error {
	// Find all descendants
	filter := bson.M{
		"path":       bson.M{"$regex": "^" + oldPath + "/"},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var descendants []*models.ProductCategory
	if err = cursor.All(ctx, &descendants); err != nil {
		return err
	}

	// Update each descendant's path
	for _, desc := range descendants {
		newDescPath := strings.Replace(desc.Path, oldPath, newPath, 1)
		_, err := r.collection.UpdateOne(
			ctx,
			bson.M{"_id": desc.ID},
			bson.M{
				"$set": bson.M{
					"path":       newDescPath,
					"updated_at": time.Now(),
				},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete soft deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()

	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// HasChildren checks if a category has any children
func (r *CategoryRepository) HasChildren(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"parent_id":  id,
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasProducts checks if a category has any products
func (r *CategoryRepository) HasProducts(ctx context.Context, categoryID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"category_id": categoryID,
		"deleted_at":  nil,
	}

	count, err := r.productCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasActiveProducts checks if a category has any active products
func (r *CategoryRepository) HasActiveProducts(ctx context.Context, categoryID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"category_id": categoryID,
		"is_active":   true,
		"deleted_at":  nil,
	}

	count, err := r.productCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateProductCount updates the product count for a category
func (r *CategoryRepository) UpdateProductCount(ctx context.Context, categoryID primitive.ObjectID) error {
	filter := bson.M{
		"category_id": categoryID,
		"deleted_at":  nil,
	}

	count, err := r.productCollection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"product_count": count,
			"updated_at":    time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": categoryID}, update)
	return err
}

// GetCategoryTree retrieves the full category tree for an organization
func (r *CategoryRepository) GetCategoryTree(ctx context.Context, orgID primitive.ObjectID) ([]*models.ProductCategory, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	opts := options.Find().SetSort(bson.D{{Key: "level", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// CheckNameExists checks if a category name exists in the organization
func (r *CategoryRepository) CheckNameExists(ctx context.Context, orgID primitive.ObjectID, name string, parentID *primitive.ObjectID, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"name":            name,
		"deleted_at":      nil,
	}

	if parentID != nil {
		filter["parent_id"] = *parentID
	} else {
		filter["parent_id"] = nil
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindByIDs retrieves multiple categories by their IDs
func (r *CategoryRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.ProductCategory, error) {
	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// HasProductsInCategories checks if any of the given categories have products
func (r *CategoryRepository) HasProductsInCategories(ctx context.Context, categoryIDs []primitive.ObjectID) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, categoryID := range categoryIDs {
		hasProducts, err := r.HasProducts(ctx, categoryID)
		if err != nil {
			return nil, err
		}
		result[categoryID.Hex()] = hasProducts
	}

	return result, nil
}

// DeleteMultiple soft deletes multiple categories
func (r *CategoryRepository) DeleteMultiple(ctx context.Context, ids []primitive.ObjectID) error {
	now := time.Now()

	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}
