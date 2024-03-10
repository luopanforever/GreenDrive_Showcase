package service

import (
	"io"

	"github.com/luopanforever/backgreendrive/model"
	"github.com/luopanforever/backgreendrive/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CarService struct {
	CarRepo *repository.CarRepository
}

// NewCarService creates a new car service.
func NewCarService() *CarService {
	carRepo := repository.NewCarRepository()
	return &CarService{CarRepo: carRepo}
}

// GetCarModelByID gets a car model by ID.
func (s *CarService) GetCarModelByID(id primitive.ObjectID) (model.CarMetadata, io.Reader, error) {
	// id, err := primitive.ObjectIDFromHex(idStr)
	// if err != nil {
	// 	return model.CarMetadata{}, nil, err
	// }

	return s.CarRepo.FindCarModelByID(id)
}

// GetCarIdByFileName gets a car's ID by its file name.
func (s *CarService) GetCarIdByFileName(fileName string) (primitive.ObjectID, error) {
	return s.CarRepo.FindCarIdByFileName(fileName)
}
