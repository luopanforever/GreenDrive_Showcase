package service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/luopanforever/backgreendrive/common"
	"github.com/luopanforever/backgreendrive/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DownloadService struct {
	ModelRepo *repository.ModelRepository
	ShowRepo  *repository.ShowRepository
	AO        *common.AliOss
}

func NewDownloadService() *DownloadService {
	modelRepo := repository.GetModelRepository()
	showRepo := repository.GetShowRepository()
	aO := common.GetAliOss()
	return &DownloadService{
		ModelRepo: modelRepo,
		ShowRepo:  showRepo,
		AO:        aO,
	}
}

func (s *DownloadService) DownloadModelAndResources(carName string) (string, error) {
	modelData, err := s.ModelRepo.FindModelDataByCarName(carName + ".gltf")
	if err != nil {
		return "", err
	}

	// 清除alioss上objectname为carname的zip文件
	err = s.AO.DeleteAliOss(carName)
	if err != nil {
		return "", err
	}

	// 创建临时目录用于保存文件
	tempDir := filepath.Join("/tmp/car/download", carName)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}

	// 从GridFS下载文件并保存到临时目录
	idList := s.ModelRepo.GetIdListByModelData(*modelData)
	err = s.saveFilesToTempDir(idList, tempDir)
	if err != nil {
		return "", err
	}

	// 将临时目录中的文件添加到zip文件中
	zipFilePath := filepath.Join("/tmp/car/download", fmt.Sprintf("%s.zip", carName))
	err = s.createZipFromTempDir(tempDir, zipFilePath)
	if err != nil {
		return "", err
	}

	// 将zip上传到alioss上
	fileUri, err := s.AO.AddAliOss(carName, zipFilePath)
	if err != nil {
		return "", err
	}

	// 使用dwhere模块的新函数来请求模型转换并查询结果
	outfile, err := common.ConvertAndQueryModel(fileUri, "obj") // 假设输出格式是 "obj"
	if err != nil {
		return "", fmt.Errorf("failed to convert and query model: %v", err)
	}

	return outfile, nil
}

// 该函数可能位于DownloadService中
func (s *DownloadService) saveFilesToTempDir(idList []primitive.ObjectID, tempDir string) error {
	for _, id := range idList {
		carMeta, dStream, err := s.ShowRepo.FindCarModelByID(id)
		if err != nil {
			return err
		}

		// 处理文件名，对于特殊文件名进行重命名
		filename := carMeta.Filename
		if strings.HasSuffix(filename, ".gltf") {
			filename = "scene.gltf" // 将car?.gltf重命名为scene.gltf
		} else if strings.HasPrefix(filename, "textures/") {
			// 确保textures目录存在
			texturesDir := filepath.Join(tempDir, "textures")
			if err := os.MkdirAll(texturesDir, 0755); err != nil {
				return err
			}
			filename = filepath.Join("textures", filepath.Base(filename)) // 保留textures/前缀
		}

		// 保存文件到临时目录
		filePath := filepath.Join(tempDir, filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, dStream); err != nil {
			return err
		}
	}

	return nil
}

// 该函数可能位于DownloadService中
func (s *DownloadService) createZipFromTempDir(tempDir, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历临时目录中的所有文件和目录，将它们添加到zip中
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过根目录
		if path == tempDir {
			return nil
		}

		// 创建zip文件中的文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置zip中文件的相对路径
		header.Name, err = filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/" // 明确指定这是一个目录
		} else {
			header.Method = zip.Deflate // 使用Deflate压缩算法压缩文件
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, _ = io.Copy(writer, file)
		}

		return err
	})

	return err
}
