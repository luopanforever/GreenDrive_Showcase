package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/config"
)

func main() {
	r := gin.Default()
	// r.Use(common.CORSMiddleware())

	// 初始化mongodb连接
	config.ConnectDB()
	defer func() {
		if err := config.MongoDB.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error on disconnecting from MongoDB: %v", err)
		}
	}()
	// 获取MongoDB数据库实例

	db := config.MongoDB.Database("tdCars")
	r = CollectRoute(r, db)

	r.Run() // 默认在0.0.0.0:8080启动服务
}
