package controller

import (
	"net/http"
	"strings"

	"github.com/luopanforever/backgreendrive/response"
	"github.com/luopanforever/backgreendrive/service"

	"github.com/gin-gonic/gin"
)

type FileChunkController struct {
	CarService *service.FileChunkService
}

// NewCarController creates a new car controller.
func NewShowController() *FileChunkController {
	carService := service.NewShowService()
	return &FileChunkController{CarService: carService}
}

// GetCarModel handles the request to get a car model by ID.
func (cc *FileChunkController) GetCarModelByFileName(c *gin.Context) {
	// _ = c.Param("carId")
	action := c.Param("action")
	fileName := strings.TrimPrefix(action, "/")
	// println(fileName)

	carId, err := cc.CarService.GetCarIdByFileName(fileName)
	if err != nil {
		response.Fail(c, "Failed to get the car id", gin.H{"error": err.Error()})
		return
	}

	carMeta, file, err := cc.CarService.GetCarModelByID(carId)
	if err != nil {
		response.Fail(c, "Failed to GetCarModelByID", gin.H{"error": err.Error()})
		return
	}

	// Streaming the file to the client
	c.DataFromReader(http.StatusOK, carMeta.Length, "application/octet-stream", file, nil)
}
