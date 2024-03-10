package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
	"github.com/luopanforever/backgreendrive/repository"
	"github.com/luopanforever/backgreendrive/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func CollectRoute(r *gin.Engine, db *mongo.Database) *gin.Engine {
	carRepo := repository.NewCarRepository(db)
	carService := service.NewCarService(carRepo)
	carController := controller.NewCarController(carService)

	r.GET("/car/:filename", carController.GetCarModelByFileName)

	// r.GET("/api/test", controller.Api)

	// r.GET("/api/test", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "hello world",
	// 	})
	// })

	return r
}
