package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/shared/models"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.Status = models.ProductStatusActive
	product.TrackInventory = true

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *ProductRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	filter := bson.M{"_id": id, "deleted_at": nil}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	return &product, err
}

func (r *ProductRepository) FindBySKU(ctx context.Context, orgID primitive.ObjectID, sku string) (*models.Product, error) {
	var product models.Product
	filter := bson.M{
		"organization_id": orgID,
		"sku":             sku,
		"deleted_at":      nil,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	return &product, err
}

func (r *ProductRepository) FindAll(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Product, error) {
	filter := bson.M{"organization_id": orgID, "deleted_at": nil}
	
	if status, ok := filters["status"]; ok {
		filter["status"] = status
	}
	if categoryID, ok := filters["category_id"]; ok {
		filter["category_id"] = categoryID
	}
	if search, ok := filters["search"]; ok {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"sku": bson.M{"$regex": search, "$options": "i"}},
			{"barcode": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()
	product.Version++

	filter := bson.M{"_id": product.ID, "version": product.Version - 1}
	update := bson.M{"$set": product}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ProductRepository) UpdateStockSummary(ctx context.Context, productID primitive.ObjectID, totalStock, availableStock, allocatedStock, stockValue float64) error {
	filter := bson.M{"_id": productID}
	update := bson.M{
		"$set": bson.M{
			"total_stock":     totalStock,
			"available_stock": availableStock,
			"allocated_stock": allocatedStock,
			"stock_value":     stockValue,
			"updated_at":      time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ProductRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
			"status":     models.ProductStatusDiscontinued,
		},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ProductRepository) SKUExists(ctx context.Context, orgID primitive.ObjectID, sku string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"sku":             sku,
		"deleted_at":      nil,
	}
	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}
