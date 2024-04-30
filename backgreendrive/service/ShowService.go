package service

import (
	"io"

	"github.com/luopanforever/backgreendrive/entity"
	"github.com/luopanforever/backgreendrive/repository"
)

type ShowService struct {
	CarRepo *repository.ShowRepository
}

// NewCarService creates a new car service.
func NewShowService() *ShowService {
	carRepo := repository.GetShowRepository()
	return &ShowService{CarRepo: carRepo}
}

func (s *ShowService) GetCarModelByCarNameAndAction(carName, fileName string) (entity.CarMetadata, io.Reader, error) {
	return s.CarRepo.FindCarModelByCarNameAndAction(carName, fileName)
}
