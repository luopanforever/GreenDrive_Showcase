package controller

import (
	"net/http"
	"strings"

	"github.com/luopanforever/backgreendrive/response"
	"github.com/luopanforever/backgreendrive/service"

	"github.com/gin-gonic/gin"
)

type ShowController struct {
	CarService *service.ShowService
}

// NewCarController creates a new car controller.
func NewShowController() *ShowController {
	carService := service.NewShowService()
	return &ShowController{CarService: carService}
}

func (cc *ShowController) GetCarModelByFileName(c *gin.Context) {
	carName := c.Param("carName")
	action := c.Param("action")
	fileName := strings.TrimPrefix(action, "/")

	carMeta, file, err := cc.CarService.GetCarModelByCarNameAndAction(carName, fileName)
	if err != nil {
		response.Fail(c, "Failed to get the car model", gin.H{"error": err.Error()})
		return
	}

	c.DataFromReader(http.StatusOK, carMeta.Length, "application/octet-stream", file, nil)
}
