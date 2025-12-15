package repository

import "go.mongodb.org/mongo-driver/mongo"

type ProductRepository struct {
	collection *mongo.Collection 
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}