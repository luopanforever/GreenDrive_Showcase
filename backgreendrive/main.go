package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/config"
)

func main() {
	r := gin.Default()
	// r.Use(common.CORSMiddleware())

	// 初始化mongodb连接
	deClose := config.InItDB()
	defer config.Close(deClose)

	r = CollectRoute(r)

	r.Run() // 默认在0.0.0.0:8080启动服务
}
