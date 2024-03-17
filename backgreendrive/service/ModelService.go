package service

import (
	"github.com/luopanforever/backgreendrive/entity"
	"github.com/luopanforever/backgreendrive/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelService struct {
	Repo *repository.ModelRepository
}

func NewModelService() *ModelService {
	newFindRepo := repository.GetModelRepository()
	return &ModelService{Repo: newFindRepo}
}

// AddResourceToModel 调用ModelRepository的AddResourceToModel方法添加资源到模型
func (s *ModelService) AddResourceToModel(modelName string, resourceName string, resourceFileId primitive.ObjectID) error {
	return s.Repo.AddResourceToModel(modelName, resourceName, resourceFileId)
}

// RemoveResourceFromModel 调用ModelRepository的RemoveResourceFromModel方法从模型中移除资源
func (s *ModelService) RemoveResourceFromModel(modelName string, resourceName string) error {
	return s.Repo.RemoveResourceFromModel(modelName, resourceName)
}

// CreateModelData 调用ModelRepository的CreateModelData方法创建模型数据
func (s *ModelService) CreateModelData(modelName string, modelFileId primitive.ObjectID) error {
	return s.Repo.CreateModelData(modelName, modelFileId)
}

// DeleteModelData 调用ModelRepository的DeleteModelData方法删除模型数据
func (s *ModelService) DeleteModelData(modelName string) error {
	return s.Repo.DeleteModelData(modelName)
}

// FindModelDataByCarName 调用ModelRepository的FindModelDataByCarName方法根据汽车名称查找模型数据
func (s *ModelService) FindModelDataByCarName(carName string) (*entity.ModelData, error) {
	return s.Repo.FindModelDataByCarName(carName)
}
