package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	// 展示管理
	showController := controller.NewShowController()
	r.GET("/car/show/:carName/*action", showController.GetCarModelByFileName)

	// 上传文件
	// r.POST("/car/upload", controller.UploadController_) //后期需要优化为"/car/upload/:carId"
	uploadController := controller.NewUploadController()
	r.POST("/car/upload/:carId", uploadController.UploadZips)
	// 删除汽车所有资源
	r.DELETE("car/upload/delete/:carName", uploadController.DeleteCar)
	// 无脑删除所有filechunks资源
	// r.DELETE("/car/upload/deleteAll", uploadController.DeleteAllFiles)

	// 汽车名字管理
	nameController := controller.NewNameController()
	r.GET("/car/names/available", nameController.FindAvailableName)
	r.GET("/car/names/list", nameController.GetNameList)

	// 测试小功能点
	// carname管理测试
	// r.POST("/car/names/add/:carName", repository.GetNameRepository().AddCarNameTest)
	// r.DELETE("/car/names/remove/:carName", repository.GetNameRepository().RemoveCarNameTest)
	// modeldata管理测试
	// r.POST("/car/model/add/resource/:carName", repository.GetModelRepository().AddResourceToModelTest)
	// r.DELETE("/car/model/delete/resource/:carName/*action", repository.GetModelRepository().RemoveResourceFromModelTest)
	// r.POST("/car/model/add", repository.GetModelRepository().CreateModelDataTest)
	// r.DELETE("/car/model/delete/:modelName", repository.GetModelRepository().DeleteModelDataTest)

	return r
}
