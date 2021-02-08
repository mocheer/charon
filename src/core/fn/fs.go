package fn

import "os"

// MkdirNotExist 创建不存在的文件夹
func MkdirNotExist(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filePath, os.ModePerm)
		}
	}
	return err
}
