package service

import (
	"io"

	"github.com/luopanforever/backgreendrive/entity"
	"github.com/luopanforever/backgreendrive/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShowService struct {
	CarRepo *repository.ShowRepository
}

// NewCarService creates a new car service.
func NewShowService() *ShowService {
	carRepo := repository.GetShowRepository()
	return &ShowService{CarRepo: carRepo}
}

// GetCarModelByID gets a car model by ID.
func (s *ShowService) GetCarModelByID(id primitive.ObjectID) (entity.CarMetadata, io.Reader, error) {
	// id, err := primitive.ObjectIDFromHex(idStr)
	// if err != nil {
	// 	return model.CarMetadata{}, nil, err
	// }

	return s.CarRepo.FindCarModelByID(id)
}

// GetCarIdByFileName gets a car's ID by its file name.
func (s *ShowService) GetCarIdByFileName(fileName string) (primitive.ObjectID, error) {
	return s.CarRepo.FindCarIdByFileName(fileName)
}

func (s *ShowService) GetCarModelByCarNameAndAction(carName, fileName string) (entity.CarMetadata, io.Reader, error) {
	return s.CarRepo.FindCarModelByCarNameAndAction(carName, fileName)
}
