package repository

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/luopanforever/backgreendrive/entity" // 替换为实际的模块路径

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type ShowRepository struct {
	DB *mongo.Database
}

// NewShowRepository creates a new repository for cars.
func newShowRepository() *ShowRepository {
	return (*ShowRepository)(NewRepository())
}

func GetShowRepository() *ShowRepository {
	return newShowRepository()
}

// 通过_id查找相应资源
func (r *ShowRepository) FindCarModelByID(id primitive.ObjectID) (entity.CarMetadata, io.Reader, error) {
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

func (r *ShowRepository) FindCarModelByCarNameAndAction(carName, fileName string) (entity.CarMetadata, io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	println("进入repository中")
	bucket, err := gridfs.NewBucket(r.DB)
	if err != nil {
		return entity.CarMetadata{}, nil, err
	}
	var modelData entity.ModelData

	println("准备获取modeldata数据")
	err = r.DB.Collection("modelData").FindOne(ctx, bson.M{"modelName": carName + ".gltf"}).Decode(&modelData)
	if err != nil {
		return entity.CarMetadata{}, nil, err
	}
	// println("modeldata:")
	// println("modelname:", modelData.ModelName)
	// println("modelfileid:", modelData.ModelFileId.String())
	// println("")
	// println("开始要查找的resource文件是否存在:", fileName)
	var fileId primitive.ObjectID
	if strings.HasSuffix(fileName, ".gltf") {
		if modelData.ModelName == fileName {
			fileId = modelData.ModelFileId
			// println("匹配到了是.gltf文件")
		} else {
			return entity.CarMetadata{}, nil, err
		}
	} else {
		// println("不是.gltf文件,开始遍历每一个resource文件名")
		for _, resource := range modelData.Resources {
			println("resource name: ", resource.Name)
			if resource.Name == fileName {
				fileId = resource.FileId
				// println("匹配到了是resource文件,文件名为: ", resource.Name)
				break
			}
		}

	}
	// println("开始获取要查找资源的元信息")
	var carMeta entity.CarMetadata
	err = r.DB.Collection("fs.files").FindOne(ctx, bson.M{"_id": fileId}).Decode(&carMeta)
	if err != nil {
		return entity.CarMetadata{}, nil, err
	}
	// println("获取完毕")
	// println("开始用gridfs传输资源")
	dStream, err := bucket.OpenDownloadStream(fileId)
	if err != nil {
		return entity.CarMetadata{}, nil, err
	}
	// println("传输完毕")
	return carMeta, dStream, nil
}

// 通过汽车名查找汽车id
func (r *ShowRepository) FindCarIdByFileName(fileName string) (primitive.ObjectID, error) {
	var car entity.CarMetadata
	filter := bson.M{"filename": fileName}
	err := r.DB.Collection("fs.files").FindOne(context.Background(), filter).Decode(&car)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return car.ID, nil
}
