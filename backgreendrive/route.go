package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
	"github.com/luopanforever/backgreendrive/repository"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	// 展示管理
	showController := controller.NewShowController()
	r.GET("/car/show/:carName/*action", showController.GetCarModelByFileName)

	// 上传文件
	r.POST("/car/upload", controller.UploadController) //后期需要优化为"/car/upload/:carId"

	// 汽车名字管理
	nameController := controller.NewNameController()
	r.GET("/car/names/available", nameController.FindAvailableName)
	r.GET("/car/names/list", nameController.GetNameList)

	// 测试小功能点
	// carname管理测试
	r.POST("/car/names/add/:carName", repository.GetCarRepository().AddCarName)
	r.DELETE("/car/names/remove/:carName", repository.GetCarRepository().RemoveCarName)
	// modeldata管理测试
	r.POST("/car/model/add/resource/:carName", repository.GetCarRepository().AddResourceToModelTest)
	r.DELETE("/car/model/delete/resource/:carName/*action", repository.GetCarRepository().RemoveResourceFromModelTest)
	r.POST("/car/model/add", repository.GetCarRepository().CreateModelDataTest)
	r.DELETE("/car/model/delete/:modelName", repository.GetCarRepository().DeleteModelDataTest)

	return r
}
