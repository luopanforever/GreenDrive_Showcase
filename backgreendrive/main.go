package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/middlewares"
)

func main() {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.Run() // 默认在0.0.0.0:8080启动服务
}
