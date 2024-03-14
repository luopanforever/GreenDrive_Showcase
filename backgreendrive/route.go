package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
	"github.com/luopanforever/backgreendrive/repository"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	// 展示管理
	showController := controller.NewShowController()
	r.GET("/car/show/:carId/*action", showController.GetCarModelByFileName)

	// 汽车名字管理
	nameController := controller.NewNameController()
	r.GET("/car/names/available", nameController.FindAvailableName)
	r.GET("/car/names/list", nameController.GetNameList)

	// 测试小功能点
	r.POST("/car/add/:carName", repository.GetCarRepository().AddCarName)
	r.DELETE("/car/remove/:carName", repository.GetCarRepository().RemoveCarName)

	return r
}
