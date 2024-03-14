package service

import "github.com/luopanforever/backgreendrive/repository" // 替换为实际的包路径

type NameService struct {
	Repo *repository.CarRepository
}

func NewNameService() *NameService {
	newFindRepo := repository.GetCarRepository()
	return &NameService{Repo: newFindRepo}
}

func (service *NameService) FindAvailableName() (string, error) {
	return service.Repo.FindAvailableName()
}

func (service *NameService) GetNameList() ([]string, error) {
	return service.Repo.GetNameList()
}
