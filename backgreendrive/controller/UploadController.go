package controller

import (
	"fmt"
	"os"
	"path/filepath"
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
	carId := c.Param("carId")

	form, err := c.MultipartForm()
	if err != nil {
		response.Fail(c, "Failed to parse multipart form", gin.H{"error": err.Error()})
		return
	}

	files := form.File["file[]"] // 前端需要将文件字段命名为 file[] 以支持多文件

	for _, file := range files {
		// 对每个文件重复保存和解压缩的过程
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
		unzipDir = unzipDir + "/"
		// 处理scene.gltf文件
		// 遍历解压目录并上传其他文件
		gltfUploaded := false
		gltfPath := filepath.Join(unzipDir, "scene.gltf")

		if _, err := os.Stat(gltfPath); !os.IsNotExist(err) {
			fileId, err := ctrl.uploadService.UploadFsFileChunkModel(unzipDir, "scene.gltf", carId)
			if err != nil {
				response.Fail(c, "Failed to upload GLTF file", gin.H{"error": err.Error()})
				return
			}
			gltfUploaded = true

			// 创建modeldata记录
			err = ctrl.modelService.CreateModelData(carId, fileId)
			if err != nil {
				response.Fail(c, "Failed to create model data", gin.H{"error": err.Error()})
				return
			}

			// 在carname的数组中添加carid
			err = ctrl.nameService.Repo.AddCarName(carId)
			if err != nil {
				response.Fail(c, "Failed to add carid data", gin.H{"error": err.Error()})
				return
			}
		}

		err = filepath.Walk(unzipDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 获取文件相对于解压目录的相对路径
			relativePath, err := filepath.Rel(unzipDir, path)
			if err != nil {
				return err
			}

			// 忽略不需要上传的文件和目录
			if relativePath == "." || strings.Contains(relativePath, "__MACOSX") || strings.Contains(relativePath, ".DS_Store") {
				return nil
			}

			// 忽略license.txt文件
			if relativePath == "license.txt" {
				return nil
			}

			// 特殊处理scene.gltf文件，确保只上传一次
			if relativePath == "scene.gltf" {
				if gltfUploaded {
					return nil // 如果scene.gltf已经上传过，跳过
				}
				gltfUploaded = true // 标记scene.gltf为已上传
			}

			if info.IsDir() {
				return nil // 忽略目录本身，但不忽略其内容
			}

			// 上传文件，并获取上传后的文件ID
			fileId, err := ctrl.uploadService.UploadFsFileChunkModel(unzipDir, relativePath, carId)
			if err != nil {
				return fmt.Errorf("failed to upload file '%s': %v", relativePath, err)
			}
			fmt.Println("上传文件名为:", relativePath)
			fmt.Println("fs.files的_id为:", fileId.Hex())

			// 添加资源到modeldata文档
			err = ctrl.modelService.AddResourceToModel(carId+".gltf", relativePath, fileId)
			if err != nil {
				return fmt.Errorf("failed to add resource to model for file '%s': %v", relativePath, err)
			}

			return nil
		})

		if err != nil {
			response.Fail(c, "Failed to process unzipped files", gin.H{"error": err.Error()})
			return
		}

		carId, err = incrementNumberSuffix(carId)
		if err != nil {
			response.Fail(c, "Failed to increce carId", gin.H{"error": err.Error()})
			return
		}
	}

	response.Success(c, nil, "All files uploaded and extracted successfully")
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
	println("modeldata:")
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
