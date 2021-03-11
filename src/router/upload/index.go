package upload

import (
	"fmt"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/fs"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/router/auth"
)

// Use 初始化 uploadFile 路由
func Use(api fiber.Router) {
	router := api.Group("/upload")
	// 上传文件
	router.Post("/file/*", auth.GlobalProtected, auth.PermissProtectd, uploadFile)
	// 上传多个文件
	router.Post("/files/*", auth.GlobalProtected, auth.PermissProtectd, uploadFiles)
	// 上传文件夹（支持chrome）
	router.Post("/folder", auth.GlobalProtected, auth.PermissProtectd, uploadFolder)
}

// uploadFile 上传文件
func uploadFile(c *fiber.Ctx) error {
	var baseDir = c.Params("*", "data")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	return c.SaveFile(file, fmt.Sprintf("./%s/%s", baseDir, file.Filename))
}

// uploadFiles 上传多个文件
func uploadFiles(c *fiber.Ctx) error {
	var baseDir = c.Params("*", "data")
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	for _, file := range files {
		dst := path.Join(baseDir, file.Filename)
		dir := path.Dir(dst)
		fs.MkdirNotExist(dir)
		err := c.SaveFile(file, dst)
		if err != nil {
			return err
		}
	}
	global.Db.Exec("select * from pipal.update_app_lib_version()")

	return res.ResultOK(c, true)
}

// uploadFolder 上传文件夹
func uploadFolder(c *fiber.Ctx) error {
	var baseDir = "public"
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	for _, file := range files {
		fileName := file.Filename
		if strings.HasPrefix(fileName, "dist") {
			fileName = strings.Replace(fileName, "dist/", "", 1)
			dst := path.Join(baseDir, fileName)
			fs.MkdirNotExist(dst)
			err := c.SaveFile(file, dst)
			if err != nil {
				return err
			}
		} else {
			return res.ResultError(c, "上传的文件夹不符合规范", nil)
		}
	}
	global.Db.Exec("select * from pipal.update_app_lib_version()")
	return res.ResultOK(c, true)
}
