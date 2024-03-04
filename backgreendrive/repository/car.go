package repository

import (
	"context"
	"time"

	"github.com/luopanforever/backgreendrive/model" // 替换为实际的模块路径

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarRepository struct {
	Collection *mongo.Collection
}

func NewCarRepository(db *mongo.Client, dbName, colName string) *CarRepository {
	collection := db.Database(dbName).Collection(colName)
	return &CarRepository{Collection: collection}
}

// SaveCar saves a new car model into the database
func (r *CarRepository) SaveCar(car model.Car) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.Collection.InsertOne(ctx, car)
	return result, err
}

// FindAllCars retrieves all car models from the database
func (r *CarRepository) FindAllCars() ([]model.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var cars []model.Car
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var car model.Car
		cursor.Decode(&car)
		cars = append(cars, car)
	}

	return cars, nil
}
