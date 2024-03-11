package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	carController := controller.NewCarController()
	r.GET("/car/:carId/*action", carController.GetCarModelByFileName)

	return r
}
