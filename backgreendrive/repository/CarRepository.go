package repository

import (
	"context"
	"io"
	"time"

	"github.com/luopanforever/backgreendrive/model" // 替换为实际的模块路径

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type CarRepository struct {
	DB *mongo.Database
}

// NewCarRepository creates a new repository for cars.
func NewCarRepository(db *mongo.Database) *CarRepository {
	return &CarRepository{DB: db}
}

// FindCarModelByID retrieves a car model by its ID from GridFS.
func (r *CarRepository) FindCarModelByID(id primitive.ObjectID) (model.CarMetadata, io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Creating a GridFS bucket
	bucket, err := gridfs.NewBucket(r.DB)
	if err != nil {
		return model.CarMetadata{}, nil, err
	}

	// Finding the file metadata
	var carMeta model.CarMetadata
	err = r.DB.Collection("fs.files").FindOne(ctx, bson.M{"_id": id}).Decode(&carMeta)
	if err != nil {
		return model.CarMetadata{}, nil, err
	}

	// Downloading the file
	dStream, err := bucket.OpenDownloadStream(id)
	if err != nil {
		return model.CarMetadata{}, nil, err
	}

	return carMeta, dStream, nil
}

// FindCarIdByFileName finds a car's ID by its file name.
func (r *CarRepository) FindCarIdByFileName(fileName string) (primitive.ObjectID, error) {
	var car model.CarMetadata
	filter := bson.M{"filename": fileName}
	err := r.DB.Collection("fs.files").FindOne(context.Background(), filter).Decode(&car)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return car.ID, nil
}

// type CarRepository struct {
// 	Collection *mongo.Collection
// }

// func NewCarRepository(db *mongo.Client, dbName, colName string) *CarRepository {
// 	collection := db.Database(dbName).Collection(colName)
// 	return &CarRepository{Collection: collection}
// }

// // SaveCar saves a new car model into the database
// func (r *CarRepository) SaveCar(car model.Car) (*mongo.InsertOneResult, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	result, err := r.Collection.InsertOne(ctx, car)
// 	return result, err
// }

// // FindAllCars retrieves all car models from the database
// func (r *CarRepository) FindAllCars() ([]model.Car, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	var cars []model.Car
// 	cursor, err := r.Collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	for cursor.Next(ctx) {
// 		var car model.Car
// 		cursor.Decode(&car)
// 		cars = append(cars, car)
// 	}

// 	return cars, nil
// }

// func (r *CarRepository) FindCarByFilename(filename string) (*model.Car, error) {
// 	var car model.Car
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	filter := bson.M{"filename": filename}
// 	println(filename)
// 	err := r.Collection.FindOne(ctx, filter).Decode(&car)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &car, nil
// }
