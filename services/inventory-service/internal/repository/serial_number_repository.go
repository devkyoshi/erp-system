package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yourusername/erp-system/shared/models"
)

type SerialNumberRepository struct {
	collection *mongo.Collection
}

func NewSerialNumberRepository(db *mongo.Database) *SerialNumberRepository {
	return &SerialNumberRepository{
		collection: db.Collection("serial_numbers"),
	}
}

func (r *SerialNumberRepository) Create(ctx context.Context, serialNumber *models.SerialNumber) error {
	serialNumber.ID = primitive.NewObjectID()
	serialNumber.CreatedAt = time.Now()
	serialNumber.UpdatedAt = time.Now()
	serialNumber.IsAvailable = true
	serialNumber.Status = "available"

	_, err := r.collection.InsertOne(ctx, serialNumber)
	return err
}

func (r *SerialNumberRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.SerialNumber, error) {
	var serialNumber models.SerialNumber
	filter := bson.M{"_id": id, "deleted_at": nil}
	err := r.collection.FindOne(ctx, filter).Decode(&serialNumber)
	return &serialNumber, err
}

func (r *SerialNumberRepository) FindBySerialNo(ctx context.Context, orgID primitive.ObjectID, serialNo string) (*models.SerialNumber, error) {
	var serialNumber models.SerialNumber
	filter := bson.M{
		"organization_id": orgID,
		"serial_no":       serialNo,
		"deleted_at":      nil,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&serialNumber)
	return &serialNumber, err
}

func (r *SerialNumberRepository) FindByProduct(ctx context.Context, productID primitive.ObjectID, filters map[string]interface{}) ([]*models.SerialNumber, error) {
	filter := bson.M{
		"product_id": productID,
		"deleted_at": nil,
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		filter["status"] = status
	}
	if locationID, ok := filters["location_id"].(string); ok && locationID != "" {
		locID, _ := primitive.ObjectIDFromHex(locationID)
		filter["location_id"] = locID
	}
	if available, ok := filters["is_available"].(bool); ok {
		filter["is_available"] = available
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var serialNumbers []*models.SerialNumber
	if err = cursor.All(ctx, &serialNumbers); err != nil {
		return nil, err
	}
	return serialNumbers, nil
}

func (r *SerialNumberRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": updates}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *SerialNumberRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string, isAvailable bool) error {
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"is_available": isAvailable,
			"updated_at":   time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *SerialNumberRepository) Allocate(ctx context.Context, id primitive.ObjectID, customerID, salesOrderID primitive.ObjectID) error {
	filter := bson.M{"_id": id, "deleted_at": nil, "is_available": true}
	update := bson.M{
		"$set": bson.M{
			"status":         "allocated",
			"is_available":   false,
			"customer_id":    customerID,
			"sales_order_id": salesOrderID,
			"updated_at":     time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *SerialNumberRepository) MarkAsSold(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"status":       "sold",
			"is_available": false,
			"sold_date":    time.Now(),
			"updated_at":   time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *SerialNumberRepository) Delete(ctx context.Context, id primitive.ObjectID, deletedBy primitive.ObjectID) error {
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

func (r *SerialNumberRepository) SerialNoExists(ctx context.Context, orgID primitive.ObjectID, serialNo string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"serial_no":       serialNo,
		"deleted_at":      nil,
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
