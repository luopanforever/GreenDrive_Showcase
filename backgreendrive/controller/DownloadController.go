package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	carName := c.Param("carName")

	zipFilePath, err := ctrl.DownloadService.DownloadModelAndResources(carName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(zipFilePath)

	// 可选: 下载后删除zip文件
}
