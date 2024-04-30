package repository

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/luopanforever/backgreendrive/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type UploadRepository struct {
	DB *mongo.Database
}

// NewUploadRepository creates a new repository for cars.
func newUploadRepository() *UploadRepository {
	return (*UploadRepository)(NewRepository())
}

func GetUploadRepository() *UploadRepository {
	return newUploadRepository()
}

// 添加实打实的资源到gridfs中 /tmp/unzipped/car? 返回唯一_id
func (r *UploadRepository) UploadFsFileChunkModel(filePath, fileName, carName string) (primitive.ObjectID, error) {

	// 创建一个新的 GridFS bucket
	bucket, err := gridfs.NewBucket(r.DB)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// 打开GLTF文件
	file, err := os.Open(filePath + fileName)
	if err != nil {
		return primitive.NilObjectID, err
	}
	defer file.Close()

	// 读取文件内容
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// 只修改3d汽车模型中gltf格式文件的名字
	if strings.HasSuffix(fileName, ".gltf") {
		fileName = carName + ".gltf"
	}

	// 创建一个新的 GridFS 文件
	uploadStream, err := bucket.OpenUploadStream(fileName) // 设置 GLTF 文件的名称
	if err != nil {
		return primitive.NilObjectID, err
	}
	defer uploadStream.Close()

	// 将文件内容写入 GridFS
	_, err = uploadStream.Write(fileData)
	if err != nil {
		return primitive.NilObjectID, err
	}

	fmt.Printf("Write file to DB was successful. File name: %s\n", fileName)

	// 返回文件ID
	return uploadStream.FileID.(primitive.ObjectID), nil
}

func (r *UploadRepository) DeleteCarResources(modelData entity.ModelData) error {
	// 创建一个新的 GridFS bucket
	bucket, err := gridfs.NewBucket(r.DB)
	if err != nil {
		return err
	}

	// 删除模型文件
	if err := bucket.Delete(modelData.ModelFileId); err != nil {
		return err
	}

	// 删除其他资源文件
	for _, resource := range modelData.Resources {
		if err := bucket.Delete(resource.FileId); err != nil {
			return err
		}
	}

	fmt.Println("All car resources deleted successfully.")
	return nil
}

func (r *UploadRepository) DeleteFsFileById(fileId primitive.ObjectID) error {
	// 创建一个新的 GridFS bucket
	bucket, err := gridfs.NewBucket(r.DB)
	if err != nil {
		return err
	}

	// 删除文件
	err = bucket.Delete(fileId)
	if err != nil {
		return err
	}

	fmt.Printf("File with ID %s deleted successfully.\n", fileId.Hex())
	return nil
}

// 用于开发测试环境/car/upload/deleteAll
func (r *UploadRepository) DeleteAllFsFiles() error {
	// 创建一个新的 GridFS bucket
	bucket, err := gridfs.NewBucket(r.DB)
	if err != nil {
		return err
	}

	// 查找fs.files集合中的所有记录
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.DB.Collection("fs.files").Find(ctx, bson.D{{}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// 遍历所有记录
	for cursor.Next(ctx) {
		var fileDoc bson.M
		if err := cursor.Decode(&fileDoc); err != nil {
			return err
		}

		// 从记录中获取文件的_id
		fileId, ok := fileDoc["_id"].(primitive.ObjectID)
		if !ok {
			return fmt.Errorf("invalid file id format")
		}

		// 删除文件
		if err := bucket.Delete(fileId); err != nil {
			return err
		}
		fmt.Printf("File with ID %s deleted successfully.\n", fileId.Hex())
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

// 根据本地存储的解压文件的位置逐个上传文件到服务器中
func (r *UploadRepository) ProcessUploadsAndResources(unzipDir, carName string, modelRepository *ModelRepository, NameRepository *NameRepository) error {
	unzipDir = unzipDir + "/"
	println("unzipDir: ", unzipDir)
	gltfPath := filepath.Join(unzipDir, "scene.gltf")
	println("gltfPath: ", gltfPath)
	var carModelData entity.CarMetadata
	carModelData.Filename = carName
	// 先上传scene.gltf的资源
	if _, err := os.Stat(gltfPath); !os.IsNotExist(err) {
		carModelData.ID, err = r.UploadFsFileChunkModel(unzipDir, "scene.gltf", carName)
		if err != nil {
			return err
		}

		// 创建modeldata集合
		err = modelRepository.CreateModelData(carModelData)
		if err != nil {
			return err
		}

		// 在carNames中添加carName
		err = NameRepository.AddCarName(carName)
		if err != nil {
			return err
		}
	}
	// 逐个添加除了scene.gltf的其他资源文件
	return filepath.Walk(unzipDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

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

		if info.IsDir() {
			return nil // 忽略目录本身，但不忽略其内容
		}
		if relativePath != "scene.gltf" {
			// 上传资源到resource中
			fileId, err := r.UploadFsFileChunkModel(unzipDir, relativePath, carName)
			if err != nil {
				return err
			}
			var resourceInfo entity.ResourceInfo
			resourceInfo.FileId = fileId
			resourceInfo.Name = relativePath
			// 添加资源到modelData的resource数组
			return modelRepository.AddResourceToModel(carName+".gltf", resourceInfo)
		}

		return nil
	})
}
