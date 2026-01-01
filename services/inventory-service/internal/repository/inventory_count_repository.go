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

type InventoryCountRepository struct {
	collection *mongo.Collection
}

func NewInventoryCountRepository(db *mongo.Database) *InventoryCountRepository {
	return &InventoryCountRepository{
		collection: db.Collection("inventory_counts"),
	}
}

func (r *InventoryCountRepository) Create(ctx context.Context, count *models.InventoryCount) error {
	count.ID = primitive.NewObjectID()
	count.CreatedAt = time.Now()
	count.UpdatedAt = time.Now()

	if count.Status == "" {
		count.Status = "in_progress"
	}
	count.StartedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, count)
	return err
}

func (r *InventoryCountRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.InventoryCount, error) {
	var count models.InventoryCount
	filter := bson.M{"_id": id, "deleted_at": nil}
	err := r.collection.FindOne(ctx, filter).Decode(&count)
	return &count, err
}

func (r *InventoryCountRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.InventoryCount, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		filter["status"] = status
	}
	if locationID, ok := filters["location_id"].(string); ok && locationID != "" {
		locID, _ := primitive.ObjectIDFromHex(locationID)
		filter["location_id"] = locID
	}
	if countType, ok := filters["count_type"].(string); ok && countType != "" {
		filter["count_type"] = countType
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "count_date", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var counts []*models.InventoryCount
	if err = cursor.All(ctx, &counts); err != nil {
		return nil, err
	}
	return counts, nil
}

func (r *InventoryCountRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": updates}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *InventoryCountRepository) Complete(ctx context.Context, id primitive.ObjectID, completedBy primitive.ObjectID, summary map[string]interface{}) error {
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"status":              "completed",
			"completed_by":        completedBy,
			"completed_at":        time.Now(),
			"total_items_counted": summary["total_items_counted"],
			"total_variance":      summary["total_variance"],
			"variance_value":      summary["variance_value"],
			"updated_at":          time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *InventoryCountRepository) UpdateItem(ctx context.Context, countID primitive.ObjectID, item models.InventoryCountItem) error {
	filter := bson.M{
		"_id":              countID,
		"items.product_id": item.ProductID,
	}

	update := bson.M{
		"$set": bson.M{
			"items.$":    item,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// If item doesn't exist, add it
	if result.MatchedCount == 0 {
		filter = bson.M{"_id": countID}
		update = bson.M{
			"$push": bson.M{"items": item},
			"$set":  bson.M{"updated_at": time.Now()},
		}
		_, err = r.collection.UpdateOne(ctx, filter, update)
	}

	return err
}

func (r *InventoryCountRepository) Delete(ctx context.Context, id primitive.ObjectID, deletedBy primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *InventoryCountRepository) Count(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}) (int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		filter["status"] = status
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count, err
}
