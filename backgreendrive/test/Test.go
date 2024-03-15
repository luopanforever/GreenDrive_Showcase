package test

import (
	"fmt"

	"github.com/luopanforever/backgreendrive/repository"
)

// 用于测试该后端的某些函数
func Test() {
	filePath := "/tmp/car/unzipped/car3/" // 文件所在的目录路径
	fileName := "scene.gltf"              // 文件名
	carId := "car1"                       // 根据需要设置汽车ID
	objid, err := repository.GetUploadRepository().UploadFsFileChunkModel(filePath, fileName, carId)
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("_id:%s\n", objid)
}
