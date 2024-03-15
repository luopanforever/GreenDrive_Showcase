package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/config"
)

func main() {
	r := gin.Default()
	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许的源，根据需要调整
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 初始化mongodb连接
	deClose := config.InItDB()
	defer config.Close(deClose)

	r = CollectRoute(r)

	// 测试用
	// test.Test_delete_file_chunk_byId()

	r.Run() // 默认在0.0.0.0:8080启动服务
}
