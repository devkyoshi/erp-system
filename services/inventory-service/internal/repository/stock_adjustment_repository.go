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

type StockAdjustmentRepository struct {
	collection *mongo.Collection
}

func NewStockAdjustmentRepository(db *mongo.Database) *StockAdjustmentRepository {
	return &StockAdjustmentRepository{
		collection: db.Collection("stock_adjustments"),
	}
}

func (r *StockAdjustmentRepository) Create(ctx context.Context, adjustment *models.StockAdjustment) error {
	adjustment.ID = primitive.NewObjectID()
	adjustment.CreatedAt = time.Now()
	adjustment.UpdatedAt = time.Now()

	if adjustment.Status == "" {
		adjustment.Status = "draft"
	}

	_, err := r.collection.InsertOne(ctx, adjustment)
	return err
}

func (r *StockAdjustmentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.StockAdjustment, error) {
	var adjustment models.StockAdjustment
	filter := bson.M{"_id": id, "deleted_at": nil}
	err := r.collection.FindOne(ctx, filter).Decode(&adjustment)
	return &adjustment, err
}

func (r *StockAdjustmentRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.StockAdjustment, error) {
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

	opts := options.Find().
		SetSort(bson.D{{Key: "adjustment_date", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var adjustments []*models.StockAdjustment
	if err = cursor.All(ctx, &adjustments); err != nil {
		return nil, err
	}
	return adjustments, nil
}

func (r *StockAdjustmentRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": updates}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *StockAdjustmentRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string, approvedBy primitive.ObjectID) error {
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"approved_by": approvedBy,
			"approved_at": time.Now(),
			"updated_at":  time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *StockAdjustmentRepository) Delete(ctx context.Context, id primitive.ObjectID, deletedBy primitive.ObjectID) error {
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

func (r *StockAdjustmentRepository) Count(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}) (int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		filter["status"] = status
	}

	return r.collection.CountDocuments(ctx, filter)
}
