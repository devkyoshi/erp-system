package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UnitChartRepository struct {
	collection *mongo.Collection
}

func NewUnitChartRepository(db *mongo.Database) *UnitChartRepository {
	return &UnitChartRepository{
		collection: db.Collection("unit_charts"),
	}
}

func (r *UnitChartRepository) Create(ctx context.Context, c *models.UnitChart) error {
	c.ID = primitive.NewObjectID()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Version = 1
	c.IsActive = true

	_, err := r.collection.InsertOne(ctx, c)
	return err
}

func (r *UnitChartRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.UnitChart, error) {
	var chart models.UnitChart
	filter := bson.M{"_id": id, "deleted_at": nil}

	err := r.collection.FindOne(ctx, filter).Decode(&chart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("unit chart not found")
		}
		return nil, err
	}

	return &chart, nil
}

func (r *UnitChartRepository) Find(ctx context.Context, activeOnly bool) ([]*models.UnitChart, error) {
	filter := bson.M{"deleted_at": nil}
	if activeOnly {
		filter["is_active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var charts []*models.UnitChart
	if err := cursor.All(ctx, &charts); err != nil {
		return nil, err
	}

	return charts, nil
}

func (r *UnitChartRepository) FindByUnits(
	ctx context.Context,
	fromUnitID, toUnitID primitive.ObjectID,
) (*models.UnitChart, error) {
	var chart models.UnitChart
	filter := bson.M{
		"from_unit_id": fromUnitID,
		"to_unit_id":   toUnitID,
		"deleted_at":   nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&chart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("unit chart not found")
		}
		return nil, err
	}

	return &chart, nil
}

func (r *UnitChartRepository) Update(ctx context.Context, chart *models.UnitChart) error {
	chart.UpdatedAt = time.Now()
	chart.Version++

	filter := bson.M{"_id": chart.ID, "deleted_at": nil}
	update := bson.M{"$set": chart}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("unit chart not found")
	}

	return nil
}

func (r *UnitChartRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("unit chart not found")
	}

	return nil
}

func (r *UnitChartRepository) HasConversions(ctx context.Context, unitID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_unit_id": unitID},
			{"to_unit_id": unitID},
		},
		"is_active":  true,
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Used to prevent circular paths
func (r *UnitChartRepository) PathExists(
	ctx context.Context,
	from primitive.ObjectID,
	to primitive.ObjectID,
) (bool, error) {

	filter := bson.M{"from_unit_id": from, "is_active": true, "deleted_at": nil}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var chart models.UnitChart
		cursor.Decode(&chart)

		if chart.ToUnitID == to {
			return true, nil
		}

		exists, _ := r.PathExists(ctx, chart.ToUnitID, to)
		if exists {
			return true, nil
		}
	}
	return false, nil
}
