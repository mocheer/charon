package fs

import (
	"image"
	"os"
)

// IsExist 检查文件或目录是否存在
func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// CreateOrExist 判断文件是否存在，不存在则创建
func CreateOrExist(path string) {
	isExit := IsExist(path)
	if !isExit {
		os.Mkdir(path, os.ModePerm)
	}
}

// GetImageFromFilePath 读取图片返回image对象
func GetImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(f)
	return image, err
}
