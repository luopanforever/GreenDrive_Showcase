package repository

import (
	"sync"

	"github.com/luopanforever/backgreendrive/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB *mongo.Database
}

var instance *Repository
var once sync.Once

func NewRepository() *Repository {
	once.Do(func() {
		instance = &Repository{DB: config.GetDB().Database("tdCars")}
	})
	return instance
}

func GetCarRepository() *Repository {
	return NewRepository()
}
