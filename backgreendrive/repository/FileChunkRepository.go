package repository

import (
	"context"
	"io"
	"time"

	"github.com/luopanforever/backgreendrive/entity" // 替换为实际的模块路径

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type FileChunkRepository struct {
	DB *mongo.Database
}

// NewFileChunkRepository creates a new repository for cars.
func newFileChunkRepository() *FileChunkRepository {
	return (*FileChunkRepository)(NewRepository())
}

func GetFileChunkRepository() *FileChunkRepository {
	return newFileChunkRepository()
}

// 通过_id查找相应资源
func (r *FileChunkRepository) FindCarModelByID(id primitive.ObjectID) (entity.CarMetadata, io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Creating a GridFS bucket
	bucket, err := gridfs.NewBucket(r.DB)
	if err != nil {
		return entity.CarMetadata{}, nil, err
	}

	// Finding the file metadata
	var carMeta entity.CarMetadata
	err = r.DB.Collection("fs.files").FindOne(ctx, bson.M{"_id": id}).Decode(&carMeta)
	if err != nil {
		return entity.CarMetadata{}, nil, err
	}

	// Downloading the file
	dStream, err := bucket.OpenDownloadStream(id)
	if err != nil {
		return entity.CarMetadata{}, nil, err
	}

	return carMeta, dStream, nil
}

// 通过汽车名查找汽车id
func (r *FileChunkRepository) FindCarIdByFileName(fileName string) (primitive.ObjectID, error) {
	var car entity.CarMetadata
	filter := bson.M{"filename": fileName}
	err := r.DB.Collection("fs.files").FindOne(context.Background(), filter).Decode(&car)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return car.ID, nil
}
