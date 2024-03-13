package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	showController := controller.NewShowController()
	findNameController := controller.NewNameController()
	r.GET("/car/show/:carId/*action", showController.GetCarModelByFileName)
	r.GET("/car/names/available", findNameController.GetAvailableName)
	return r
}
