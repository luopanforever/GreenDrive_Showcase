package common

import (
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliOss struct {
	bucket  *oss.Bucket
	fileUri string
}

func NewAliOss() (*AliOss, error) {

	bucketName := "tdcarszips"
	fileUri := "https://tdcarszips.oss-cn-beijing.aliyuncs.com/"
	endPoint := "https://oss-cn-beijing.aliyuncs.com"
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		return nil, err
	}
	client, err := oss.New(endPoint, "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return &AliOss{bucket: bucket, fileUri: fileUri}, nil
}

func GetAliOss() *AliOss {
	bucket, err := NewAliOss()
	if err != nil {
		log.Default().Fatal(err)
		return nil
	}
	return bucket
}

func (ao *AliOss) AddAliOss(objectName, localFileName string) (string, error) {
	// objectName := "car1"
	// localFileName := "/tmp/car/download/car1.zip"
	err := ao.bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return "", err
	}
	fileUri := ao.fileUri + objectName
	return fileUri, nil
}

func (ao *AliOss) DeleteAliOss(objectName string) error {
	err := ao.bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}
