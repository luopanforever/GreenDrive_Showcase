package service

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/luopanforever/backgreendrive/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadService struct {
	UploadRepo *repository.UploadRepository
}

func NewUploadService() *UploadService {
	uplodaRepo := repository.GetUploadRepository()
	return &UploadService{UploadRepo: uplodaRepo}
}

// SaveZipFile 保存上传的ZIP文件
func (s *UploadService) SaveZipFile(fileHeader *multipart.FileHeader, carId string) (string, error) {
	// 确保存放ZIP文件的目录存在
	zipDir := "/tmp/car/zip/"
	os.MkdirAll(zipDir, os.ModePerm)

	fileName := fileHeader.Filename

	// 构造ZIP文件的存储路径
	zipFilePath := filepath.Join(zipDir, fileName)

	// 保存ZIP文件
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	outFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", err
	}

	return zipFilePath, nil
}

// UnzipFiles 解压ZIP文件到指定目录
func (s *UploadService) UnzipFiles(zipFilePath, carId string) (string, error) {
	unzipDir := fmt.Sprintf("/tmp/car/unzipped/%s", carId)
	os.MkdirAll(unzipDir, os.ModePerm)

	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// 跳过顶层目录的标志
	skipTopLevelFolder := true
	var topFolderName string

	for _, f := range r.File {
		if skipTopLevelFolder {
			// 分割路径，找到顶层目录名
			parts := strings.SplitN(f.Name, "/", 2)
			if len(parts) > 0 {
				topFolderName = parts[0] + "/"
			}
			skipTopLevelFolder = false
		}

		// 移除顶层目录部分
		innerPath := strings.TrimPrefix(f.Name, topFolderName)

		fpath := filepath.Join(unzipDir, innerPath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return "", err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}

			rc, err := f.Open()
			if err != nil {
				outFile.Close()
				return "", err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return "", err
			}
		}
	}

	return unzipDir, nil
}
func (s *UploadService) UploadFsFileChunkModel(filePath, fileName, carId string) (primitive.ObjectID, error) {
	return_id, err := s.UploadRepo.UploadFsFileChunkModel(filePath, fileName, carId)
	return return_id, err
}

func (s *UploadService) DeleteAllFiles() error {
	err := s.UploadRepo.DeleteAllFsFiles()
	return err
}
