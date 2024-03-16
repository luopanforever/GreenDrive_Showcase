package repository

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/entity"
	"github.com/luopanforever/backgreendrive/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelRepository struct {
	DB *mongo.Database
}

func newModelRepository() *ModelRepository {
	return (*ModelRepository)(NewRepository())
}

func GetModelRepository() *ModelRepository {
	return newModelRepository()
}

// model管理
func (r *ModelRepository) AddResourceToModelTest(c *gin.Context) {
	var resourceInfo entity.ResourceInfo
	if err := c.BindJSON(&resourceInfo); err != nil {
		response.Fail(c, "Failed to parse request body", gin.H{"error": err.Error()})
		return
	}
	modelName := c.Param("carName")
	resourceName := resourceInfo.Name
	resourceFileId := resourceInfo.FileId

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
func (r *ModelRepository) AddResourceToModel(modelName string, resourceName string, resourceFileId primitive.ObjectID) error {
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
	return err
}

func (r *ModelRepository) RemoveResourceFromModelTest(c *gin.Context) {
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
func (r *ModelRepository) RemoveResourceFromModel(modelName string, resourceName string) error {
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
	return err
}

// 根据modeldata和modelname更新
func (r *ModelRepository) CreateModelDataTest(c *gin.Context) {

	var resourceInfo entity.ResourceInfo
	if err := c.BindJSON(&resourceInfo); err != nil {
		response.Fail(c, "Failed to parse request body", gin.H{"error": err.Error()})
		return
	}

	modelName := resourceInfo.Name
	modelFileId := resourceInfo.FileId

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
func (r *ModelRepository) CreateModelData(modelName string, modelFileId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newModelData := bson.M{
		"modelName":   modelName + ".gltf",
		"modelFileId": modelFileId,
		"resources":   []bson.M{}, // 初始化为空的数组
	}

	_, err := r.DB.Collection("modelData").InsertOne(ctx, newModelData)
	return err
}

func (r *ModelRepository) DeleteModelDataTest(c *gin.Context) {
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
func (r *ModelRepository) DeleteModelData(modelName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"modelName": modelName}
	_, err := r.DB.Collection("modelData").DeleteOne(ctx, filter)
	return err
}

func (r *ModelRepository) FindModelDataByCarName(carName string) (*entity.ModelData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var modelData entity.ModelData
	// 注意，这里我们假设carName已经包含了.gltf扩展名
	err := r.DB.Collection("modelData").FindOne(ctx, bson.M{"modelName": carName}).Decode(&modelData)
	if err != nil {
		return nil, err
	}

	return &modelData, nil
}
