package repository

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/luopanforever/backgreendrive/config"
	"github.com/luopanforever/backgreendrive/model" // 替换为实际的模块路径

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type CarRepository struct {
	DB *mongo.Database
}

var instance *CarRepository
var once sync.Once

// NewCarRepository creates a new repository for cars.
func NewCarRepository() *CarRepository {
	once.Do(func() {
		instance = &CarRepository{DB: config.GetDB().Database("tdCars")}
	})
	return instance
}

// GetCarRepository returns the singleton instance of CarRepository.
func GetCarRepository() *CarRepository {
	return NewCarRepository() // This ensures the instance is created if it doesn't exist and returns the existing one if it does.
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

func (r *CarRepository) FindAvailableName() (string, error) {
	var result struct {
		UsedNames []string `bson:"usedNames"`
	}
	if err := r.DB.Collection("carNames").FindOne(context.Background(), bson.D{}).Decode(&result); err != nil {
		return "", err
	}

	nameMap := make(map[int]bool)
	for _, name := range result.UsedNames {
		if len(name) > 3 {
			if num, err := strconv.Atoi(name[3:]); err == nil {
				nameMap[num] = true
			}
		}
	}

	for i := 1; ; i++ {
		if !nameMap[i] {
			return fmt.Sprintf("car%02d", i), nil
		}
	}
}
