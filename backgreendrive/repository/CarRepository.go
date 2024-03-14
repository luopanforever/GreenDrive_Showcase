package repository

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/config"
	"github.com/luopanforever/backgreendrive/model" // 替换为实际的模块路径
	"github.com/luopanforever/backgreendrive/response"

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

	maxNum := 0
	for _, name := range result.UsedNames {
		if len(name) > 3 {
			// 从名称中提取编号部分，并转换为整数
			if num, err := strconv.Atoi(name[3:]); err == nil {
				// 更新最大编号
				if num > maxNum {
					maxNum = num
				}
			}
		}
	}

	// 生成下一个可用的名称，即最大编号加1
	return fmt.Sprintf("car%d", maxNum+1), nil

	// [car1, car3] return car2
	// nameMap := make(map[int]bool)
	// for _, name := range result.UsedNames {
	// 	if len(name) > 3 {
	// 		if num, err := strconv.Atoi(name[3:]); err == nil {
	// 			nameMap[num] = true
	// 		}
	// 	}
	// }

	// for i := 1; ; i++ {
	// 	if !nameMap[i] {
	// 		return fmt.Sprintf("car%02d", i), nil
	// 	}
	// }
}

// name管理
// 获取汽车名字列表
func (r *CarRepository) GetNameList() ([]string, error) {
	var result struct {
		UsedNames []string `bson:"usedNames"`
	}
	if err := r.DB.Collection("carNames").FindOne(context.Background(), bson.D{}).Decode(&result); err != nil {
		return nil, err
	}
	return result.UsedNames, nil
}

// AddCarName adds a new car name to the usedNames array in carNames collection.
func (r *CarRepository) AddCarName(c *gin.Context) {
	name := c.Param("carName")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$push": bson.M{"usedNames": name}}
	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)

	if err != nil {
		response.Fail(c, "Failed to add car name", gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"carname": name}, "add success")

}

// AddCarName adds a new car name to the usedNames array in carNames collection.
// func (r *CarRepository) AddCarName(name string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	update := bson.M{"$push": bson.M{"usedNames": name}}
// 	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)
// 	return err
// }

// RemoveCarName removes a car name from the usedNames array in carNames collection.
func (r *CarRepository) RemoveCarName(c *gin.Context) {
	name := c.Param("carName")
	println(name)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$pull": bson.M{"usedNames": name}}
	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)
	if err != nil {
		response.Fail(c, "Failed to delete car name", gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"carname": name}, "delete success")
}

// func (r *CarRepository) RemoveCarName(name string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	update := bson.M{"$pull": bson.M{"usedNames": name}}
// 	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)
// 	return err
// }

// model管理
func (r *CarRepository) AddResourceToModelTest(c *gin.Context) {
	var resourceInfo model.ResourceInfo
	if err := c.BindJSON(&resourceInfo); err != nil {
		response.Fail(c, "Failed to parse request body", gin.H{"error": err.Error()})
		return
	}
	modelName := c.Param("carName")
	resourceName := resourceInfo.ResourceName
	resourceFileId := resourceInfo.FileID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$push": bson.M{
			"resources": bson.M{
				"name":   resourceName,
				"fileId": resourceFileId,
			},
		},
	}
	_, err := r.DB.Collection("modelData").UpdateOne(ctx, bson.M{"modelName": modelName}, update)
	if err != nil {
		response.Fail(c, "Failed to add model resource", gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"resource name": modelName, "objectId": resourceFileId}, "add success")
}

// func (r *CarRepository) AddResourceToModel(modelName string, resourceName string, resourceFileId primitive.ObjectID) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	update := bson.M{
// 		"$push": bson.M{
// 			"resources": bson.M{
// 				"name":   resourceName,
// 				"fileId": resourceFileId,
// 			},
// 		},
// 	}
// 	_, err := r.DB.Collection("modelData").UpdateOne(ctx, bson.M{"modelName": modelName}, update)
// 	return err
// }

func (r *CarRepository) RemoveResourceFromModelTest(c *gin.Context) {
	modelName := c.Param("carName")
	action := c.Param("action")
	resourceName := strings.TrimPrefix(action, "/")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"resources": bson.M{
				"name": resourceName,
			},
		},
	}
	_, err := r.DB.Collection("modelData").UpdateOne(ctx, bson.M{"modelName": modelName}, update)
	if err != nil {
		response.Fail(c, "Failed to delete model resource", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{"delete resource name": resourceName}, "delete success")
}

// func (r *CarRepository) RemoveResourceFromModel(modelName string, resourceName string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	update := bson.M{
// 		"$pull": bson.M{
// 			"resources": bson.M{
// 				"name": resourceName,
// 			},
// 		},
// 	}
// 	_, err := r.DB.Collection("modelData").UpdateOne(ctx, bson.M{"modelName": modelName}, update)
// 	return err
// }

func (r *CarRepository) CreateModelDataTest(c *gin.Context) {

	var resourceInfo model.ResourceInfo
	if err := c.BindJSON(&resourceInfo); err != nil {
		response.Fail(c, "Failed to parse request body", gin.H{"error": err.Error()})
		return
	}

	modelName := resourceInfo.ResourceName
	modelFileId := resourceInfo.FileID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newModelData := bson.M{
		"modelName":   modelName,
		"modelFileId": modelFileId,
		"resources":   []bson.M{}, // 初始化为空的数组
	}

	_, err := r.DB.Collection("modelData").InsertOne(ctx, newModelData)
	if err != nil {
		response.Fail(c, "Failed to create model data", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{"car name": modelName, "objectId": modelFileId}, "create success")
}

// func (r *CarRepository) CreateModelData(modelName string, modelFileId primitive.ObjectID) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     newModelData := bson.M{
//         "modelName":   modelName,
//         "modelFileId": modelFileId,
//         "resources":   []bson.M{}, // 初始化为空的数组
//     }

//     _, err := r.DB.Collection("modelData").InsertOne(ctx, newModelData)
//     return err
// }

func (r *CarRepository) DeleteModelDataTest(c *gin.Context) {
	modelName := c.Param("modelName")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"modelName": modelName}
	_, err := r.DB.Collection("modelData").DeleteOne(ctx, filter)
	if err != nil {
		response.Fail(c, "Failed to drop model data", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{"car name": modelName}, "drop success")
}

// func (r *CarRepository) DeleteModelData(modelName string) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     filter := bson.M{"modelName": modelName}
//     _, err := r.DB.Collection("modelData").DeleteOne(ctx, filter)
//     return err
// }
