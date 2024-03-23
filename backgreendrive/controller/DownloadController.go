package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/response"
	"github.com/luopanforever/backgreendrive/service"
)

type DownloadController struct {
	DownloadService *service.DownloadService
}

func NewDownloadController() *DownloadController {
	downloadController := service.NewDownloadService()
	return &DownloadController{DownloadService: downloadController}
}

func (ctrl *DownloadController) DownloadModel(c *gin.Context) {
	format := c.Param("format")
	carName := c.Param("carName")

	// 先清空临时目录
	downloadBase := "/tmp/car/download/"
	err := clearDirectory(downloadBase)
	if err != nil {
		response.Fail(c, "Failed to clear /tmp/car/download/ directory", gin.H{"error": err.Error()})
		return
	}

	// 调用 DownloadService 获取下载链接或文件路径
	result, err := ctrl.DownloadService.DownloadModelAndResources(carName, format, 180) // 第三个参数表示180秒超时
	if err != nil {
		response.Fail(c, "fail to download model and convert format", gin.H{"error": err.Error()})
		return
	}

	if format == "gltf" {
		// 直接返回zip文件
		c.File(result)
	} else {
		// 返回转换后的模型下载链接
		response.Success(c, gin.H{"fileUri": result}, "Download successful")
	}
}
