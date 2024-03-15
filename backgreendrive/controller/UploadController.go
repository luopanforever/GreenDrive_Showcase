package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/response"
	"github.com/luopanforever/backgreendrive/service"
)

type UploadController struct {
	uploadService *service.UploadService
}

func NewUploadController(uploadService *service.UploadService) *UploadController {
	return &UploadController{uploadService: uploadService}
}

func (ctrl *UploadController) UploadZip(c *gin.Context) {
	carId := c.Param("carId")
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "Failed to retrieve file", gin.H{"error": err.Error()})
		return
	}

	// 保存ZIP文件
	zipFilePath, err := ctrl.uploadService.SaveZipFile(file, carId)
	if err != nil {
		response.Fail(c, "Failed to save zip file", gin.H{"error": err.Error()})
		return
	}

	// 解压ZIP文件
	unzipDir, err := ctrl.uploadService.UnzipFiles(zipFilePath, carId)
	if err != nil {
		response.Fail(c, "Failed to unzip files", gin.H{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{"unzipDir": unzipDir}, "Files uploaded and extracted successfully")
}

func UploadController_(c *gin.Context) {
	println("进来了")
	// 解析表单数据，1 << 30 设置最大内存限制为1GB
	if err := c.Request.ParseMultipartForm(1 << 30); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error: " + err.Error()})
		return
	}

	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		response.Fail(c, "Failed to resolve multifiles", gin.H{"error": err.Error()})
		return
	}
	files := form.File["upload[]"] // "upload[]" 是前端表单中文件输入字段的名称

	for _, file := range files {
		println(file.Filename)
	}

	// for _, file := range files {
	// 	// 处理每个文件，例如保存到服务器
	// 	path := "./uploads/" + file.Filename
	// 	if err := c.SaveUploadedFile(file, path); err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Saving file error: " + err.Error()})
	// 		continue // 处理下一个文件
	// 	}

	// 	// 这里可以添加其他文件处理逻辑，如更新数据库等
	// }

	// 所有文件处理完毕

	response.Success(c, gin.H{"message": "All files uploaded successfully"}, "upload success")
}
