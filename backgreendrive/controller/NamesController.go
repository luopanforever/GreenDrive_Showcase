package controller

import (
	"sort"

	"github.com/luopanforever/backgreendrive/response" // 替换为实际的包路径
	"github.com/luopanforever/backgreendrive/service"  // 替换为实际的包路径

	"github.com/gin-gonic/gin"
)

type NameController struct {
	NameService *service.NameService
}

func NewNameController() *NameController {
	nameService := service.NewNameService()
	return &NameController{NameService: nameService}
}

func (ctrl *NameController) FindAvailableName(c *gin.Context) {
	name, err := ctrl.NameService.FindAvailableName()
	if err != nil {
		response.Fail(c, "Failed to find available name", gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"availableName": name}, "Available name found")
}

func (ctrl *NameController) GetNameList(c *gin.Context) {
	names, err := ctrl.NameService.GetNameList()
	if err != nil {
		response.Fail(c, "Failed to get name list", gin.H{"error": err.Error()})
		return
	}
	sort.Strings(names)
	response.Success(c, gin.H{"names": names}, "Name list retrieved successfully")
}
