package controller

import (
	"net/http"

	"github.com/luopanforever/backgreendrive/response"
	"github.com/luopanforever/backgreendrive/service"

	"github.com/gin-gonic/gin"
)

type CarController struct {
	CarService *service.CarService
}

// NewCarController creates a new car controller.
func NewCarController() *CarController {
	carService := service.NewCarService()
	return &CarController{CarService: carService}
}

// GetCarModel handles the request to get a car model by ID.
func (cc *CarController) GetCarModelByFileName(c *gin.Context) {
	fileName := c.Param("filename") // Assuming the route parameter is named 'id'

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
func Api(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}
