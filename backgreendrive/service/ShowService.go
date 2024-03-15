package service

import (
	"io"

	"github.com/luopanforever/backgreendrive/entity"
	"github.com/luopanforever/backgreendrive/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileChunkService struct {
	CarRepo *repository.FileChunkRepository
}

// NewCarService creates a new car service.
func NewShowService() *FileChunkService {
	carRepo := repository.GetFileChunkRepository()
	return &FileChunkService{CarRepo: carRepo}
}

// GetCarModelByID gets a car model by ID.
func (s *FileChunkService) GetCarModelByID(id primitive.ObjectID) (entity.CarMetadata, io.Reader, error) {
	// id, err := primitive.ObjectIDFromHex(idStr)
	// if err != nil {
	// 	return model.CarMetadata{}, nil, err
	// }

	return s.CarRepo.FindCarModelByID(id)
}

// GetCarIdByFileName gets a car's ID by its file name.
func (s *FileChunkService) GetCarIdByFileName(fileName string) (primitive.ObjectID, error) {
	return s.CarRepo.FindCarIdByFileName(fileName)
}
