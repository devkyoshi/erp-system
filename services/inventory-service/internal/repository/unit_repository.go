package repository

import (
	"context"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UnitRepository struct {
	collection *mongo.Collection
}

func NewUnitRepository(db *mongo.Database) *UnitRepository {
	return &UnitRepository{
		collection: db.Collection("units"),
	}
}

func (r *UnitRepository) Create(ctx context.Context, u *models.Unit) error {
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Version = 1
	u.IsActive = true

	_, err := r.collection.InsertOne(ctx, u)
	return err
}

func (r *UnitRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Unit, error) {
	var unit models.Unit
	err := r.collection.FindOne(ctx, bson.M{
		"_id": id, "is_active": true,
	}).Decode(&unit)

	return &unit, err
}

func (r *UnitRepository) ExistsBaseUnit(ctx context.Context, unitType string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"unit_type":    unitType,
		"is_base_unit": true,
		"is_active":    true,
	})
	return count > 0, err
}

func (r *UnitRepository) HasConversions(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.collection.Database().
		Collection("unit_charts").
		CountDocuments(ctx, bson.M{
			"$or": []bson.M{
				{"from_unit_id": id},
				{"to_unit_id": id},
			},
			"is_active": true,
		})
	return count > 0, err
}

func (r *UnitRepository) Update(ctx context.Context, unit *models.Unit) error {
	unit.UpdatedAt = time.Now()

	filter := bson.M{"_id": unit.ID, "is_active": true}
	update := bson.M{"$set": unit}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UnitRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	})
	return err
}

func (r *UnitRepository) FindByID(
	ctx context.Context,
	id primitive.ObjectID,
) (*models.Unit, error) {

	var unit models.Unit
	err := r.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"is_active": true,
	}).Decode(&unit)

	if err != nil {
		return nil, err
	}

	return &unit, nil
}

func (r *UnitRepository) Find(
	ctx context.Context,
	unitType *string,
	activeOnly bool,
) ([]*models.Unit, error) {

	filter := bson.M{}

	if unitType != nil {
		filter["unit_type"] = *unitType
	}

	if activeOnly {
		filter["is_active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var units []*models.Unit
	if err := cursor.All(ctx, &units); err != nil {
		return nil, err
	}

	return units, nil
}
