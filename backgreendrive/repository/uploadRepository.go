package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

// 添加实打实的资源 /tmp/unzipped/car?
func (r *UploadRepository) UploadFsFileChunkModel(filePath, fileName, carId string) (primitive.ObjectID, error) {

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

	if strings.HasSuffix(fileName, ".gltf") {
		fileName = carId + ".gltf"
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
