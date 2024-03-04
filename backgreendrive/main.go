package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/common"
	"github.com/luopanforever/backgreendrive/config"
)

func main() {
	r := gin.Default()
	r.Use(common.CORSMiddleware())

	// 初始化mongodb连接
	config.ConnectDB()
	if config.MongoDB != nil {
		log.Println("Connected to MongoDB")
	}

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.Run() // 默认在0.0.0.0:8080启动服务
}
