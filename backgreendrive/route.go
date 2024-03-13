package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
	"github.com/luopanforever/backgreendrive/repository"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	showController := controller.NewShowController()
	nameController := controller.NewNameController()
	r.GET("/car/show/:carId/*action", showController.GetCarModelByFileName)
	r.GET("/car/names/available", nameController.GetAvailableName)

	//testing
	r.POST("/car/add/:carName", repository.GetCarRepository().AddCarName)
	r.DELETE("/car/remove/:carName", repository.GetCarRepository().RemoveCarName)

	return r
}
