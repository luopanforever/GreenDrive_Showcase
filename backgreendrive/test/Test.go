package test

import (
	"fmt"
	"log"

	"github.com/luopanforever/backgreendrive/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 用于测试该后端的某些函数
func Test_add_file_chunk() {
	filePath := "/tmp/car/unzipped/car3/"                // 文件所在的目录路径
	fileName := "textures/forMayaAOblinn6_occlusion.png" // 文件名
	carId := "car1"                                      // 根据需要设置汽车ID
	objid, err := repository.GetUploadRepository().UploadFsFileChunkModel(filePath, fileName, carId)
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("_id:%s\n", objid)
}
func Test_delete_file_chunk_byId(fileidstr string) {
	fileIdStr := fileidstr
	fileId, err := primitive.ObjectIDFromHex(fileIdStr)
	if err != nil {
		log.Fatalf("Invalid file ID: %v", err)
	}
	err = repository.GetUploadRepository().DeleteFsFileById(fileId)
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
	}
}

func TestRun() {
	Test_delete_file_chunk_byId("65f44c0ff14c843bb46defe0")
	Test_delete_file_chunk_byId("65f44c0ff14c843bb46defe3")
	Test_delete_file_chunk_byId("65f44c0ff14c843bb46defe5")
	Test_delete_file_chunk_byId("65f44c0ff14c843bb46defe7")

}
