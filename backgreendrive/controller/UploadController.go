package controller

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/response"
	"github.com/luopanforever/backgreendrive/service"
)

type UploadController struct {
	uploadService *service.UploadService
	modelService  *service.ModelService
	nameService   *service.NameService
}

func NewUploadController() *UploadController {
	uploadService := service.NewUploadService()
	modelService := service.NewModelService()
	nameSErvice := service.NewNameService()
	return &UploadController{
		uploadService: uploadService,
		modelService:  modelService,
		nameService:   nameSErvice,
	}
}

func (ctrl *UploadController) UploadZips(c *gin.Context) {
	// 开局清空/tmp/car/文件夹
	// 清空临时目录
	zipDir := "/tmp/car/zip/"
	unzipDirBase := "/tmp/car/unzipped/"

	err := clearDirectory(zipDir)
	if err != nil {
		response.Fail(c, "Failed to clear /tmp/car/zip/ directory", gin.H{"error": err.Error()})
		return
	}
	err = clearDirectory(unzipDirBase)
	if err != nil {
		response.Fail(c, "Failed to clear /tmp/car/unzipped/ directory", gin.H{"error": err.Error()})
		return
	}

	carId := c.Param("carId")

	// 开局检查carnames里面有没有carId
	exists, err := ctrl.nameService.CarNameExists(carId)
	if err != nil {
		response.Fail(c, "Failed to check car name existence", gin.H{"error": err.Error()})
		return
	}

	if exists {
		response.Fail(c, "Car name already exists", gin.H{"error": fmt.Errorf("car name '%s' has already been used", carId).Error()})
		return
	}

	// 开始获取资源
	form, err := c.MultipartForm()
	if err != nil {
		response.Fail(c, "Failed to parse multipart form", gin.H{"error": err.Error()})
		return
	}

	files := form.File["file[]"]
	fileName := make([]string, 0)
	for _, file := range files {
		fileName = append(fileName, file.Filename)
		zipFilePath, err := ctrl.uploadService.SaveZipFile(file, carId)
		if err != nil {
			response.Fail(c, "Failed to save zip file", gin.H{"error": err.Error()})
			return
		}

		unzipDir, err := ctrl.uploadService.UnzipFiles(zipFilePath, carId)
		if err != nil {
			response.Fail(c, "Failed to unzip files", gin.H{"error": err.Error()})
			return
		}

		err = ctrl.uploadService.ProcessUploadsAndResources(unzipDir, carId, ctrl.modelService, ctrl.nameService)
		if err != nil {
			response.Fail(c, "Failed to process uploads and resources", gin.H{"error": err.Error()})
			return
		}

		carId, err = incrementNumberSuffix(carId)
		if err != nil {
			response.Fail(c, "Failed to increment carId", gin.H{"error": err.Error()})
			return
		}
	}

	response.Success(c, gin.H{"add zips": fileName}, "All files uploaded and resources processed successfully")
}

// func UploadController_(c *gin.Context) {
// 	// 解析表单数据，1 << 30 设置最大内存限制为1GB
// 	if err := c.Request.ParseMultipartForm(1 << 30); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error: " + err.Error()})
// 		return
// 	}

// 	// 获取上传的文件
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		response.Fail(c, "Failed to resolve multifiles", gin.H{"error": err.Error()})
// 		return
// 	}
// 	files := form.File["upload[]"] // "upload[]" 是前端表单中文件输入字段的名称

// 	for _, file := range files {
// 		// 处理每个文件，例如保存到服务器
// 		path := "./uploads/" + file.Filename
// 		if err := c.SaveUploadedFile(file, path); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Saving file error: " + err.Error()})
// 			continue // 处理下一个文件
// 		}

// 		// 这里可以添加其他文件处理逻辑，如更新数据库等
// 	}

// 	// 所有文件处理完毕

// 	response.Success(c, gin.H{"message": "All files uploaded successfully"}, "upload success")
// }

func (ctrl *UploadController) DeleteAllFiles(c *gin.Context) {
	err := ctrl.uploadService.DeleteAllFiles()
	if err != nil {
		response.Fail(c, "Failed to delete all files", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, nil, "All files deleted successfully")
}

func (ctrl *UploadController) DeleteCar(c *gin.Context) {
	carName := c.Param("carName") + ".gltf"

	// 获取汽车资源的modelData
	modelData, err := ctrl.modelService.FindModelDataByCarName(carName)
	if err != nil {
		response.Fail(c, "Failed to find car model data", gin.H{"error": err.Error()})
		return
	}
	// 删除汽车所有资源
	if err := ctrl.uploadService.DeleteCarResources(*modelData); err != nil {
		response.Fail(c, "Failed to delete car resources", gin.H{"error": err.Error()})
		return
	}

	// 删除modelData记录
	if err := ctrl.modelService.DeleteModelData(carName); err != nil {
		response.Fail(c, "Failed to delete car model data", gin.H{"error": err.Error()})
		return
	}

	// 从carNames中移除该汽车名
	if err := ctrl.nameService.RemoveCarName(strings.TrimSuffix(carName, ".gltf")); err != nil {
		response.Fail(c, "Failed to remove car name", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, nil, "Car and all related resources deleted successfully")
}

func incrementNumberSuffix(str string) (string, error) {
	// 找到字符串中第一个数字的位置
	index := strings.IndexFunc(str, func(r rune) bool {
		return r >= '0' && r <= '9'
	})

	// 如果没有找到数字，返回错误
	if index == -1 {
		return "", fmt.Errorf("no number found in string")
	}

	// 提取数字部分和非数字部分
	prefix := str[:index]
	numberStr := str[index:]

	// 将数字字符串转换为整数
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return "", err
	}

	// 数字加1
	number++

	// 将结果拼接回字符串
	return fmt.Sprintf("%s%d", prefix, number), nil
}
func clearDirectory(dirPath string) error {
	// 删除目录及其包含的所有内容
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}

	// 重新创建目录
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	return nil
}
