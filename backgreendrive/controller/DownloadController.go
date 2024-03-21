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
	downloadBase := "/tmp/car/download/"

	err := clearDirectory(downloadBase)
	if err != nil {
		response.Fail(c, "Failed to clear /tmp/car/download/ directory", gin.H{"error": err.Error()})
		return
	}

	carName := c.Param("carName")

	zipFilePath, err := ctrl.DownloadService.DownloadModelAndResources(carName)
	if err != nil {
		response.Fail(c, "fail to download model", gin.H{"error": err.Error()})
		return
	}
	println("zipFilePath:", zipFilePath)
	response.Success(c, gin.H{"fileUri": zipFilePath}, "shangchuanchenggong")

}
