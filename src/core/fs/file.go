package fs

import (
	"image"
	"os"
	"path/filepath"
)

// IsExist 检查文件或目录是否存在
func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// OpenOrCreate 创建不存在的文件
func OpenOrCreate(name string, flag int, perm os.FileMode) (*os.File, error) {
	isExit := IsExist(name)
	if !isExit {
		err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
		if err == nil {
			file, err := os.Create(name)
			return file, err
		}
		return nil, err
	}
	return os.OpenFile(name, flag, perm)
}

// MkdirNotExist 创建不存在的文件夹
func MkdirNotExist(path string) error {
	isExit := IsExist(path)
	if !isExit {
		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetImageFromPath 读取图片返回image对象
func GetImageFromPath(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(f)
	return image, err
}

// SaveFile 保存图片
func SaveFile(path string, data []byte) error {
	f, err := OpenOrCreate(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err == nil {
		f.Write(data)
	}
	defer f.Close()
	return err
}
