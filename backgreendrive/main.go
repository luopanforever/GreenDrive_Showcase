package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/common"
	"github.com/luopanforever/backgreendrive/config"
	"github.com/luopanforever/backgreendrive/controller"
	"github.com/luopanforever/backgreendrive/repository"
	"github.com/luopanforever/backgreendrive/service"
)

func main() {
	r := gin.Default()
	r.Use(common.CORSMiddleware())

	// 初始化mongodb连接
	config.ConnectDB()
	defer func() {
		if err := config.MongoDB.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error on disconnecting from MongoDB: %v", err)
		}
	}()
	// 获取MongoDB数据库实例
	db := config.MongoDB.Database("tdCars")

	carRepo := repository.NewCarRepository(db)
	carService := service.NewCarService(carRepo)
	carController := controller.NewCarController(carService)

	r.GET("/car/:filename", carController.GetCarModelByFileName)

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.Run() // 默认在0.0.0.0:8080启动服务
}
