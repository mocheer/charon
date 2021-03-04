package loader

import (
	"io"
	"net/http"
	"os"
)

// LoadAndSave 下载保存
func LoadAndSave(url string, filePath string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	io.Copy(file, res.Body)
	return nil
}

// Load 下载
func Load(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body) //body 拿到请求返回的内容
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
