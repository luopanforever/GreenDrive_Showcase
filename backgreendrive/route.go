package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	carController := controller.NewCarController()
	r.GET("/car/:filename", carController.GetCarModelByFileName)

	// r.GET("/api/test", controller.Api)

	// r.GET("/api/test", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "hello world",
	// 	})
	// })

	return r
}
