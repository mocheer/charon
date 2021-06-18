package upload

import (
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/pluto/fs"
)

// Use 初始化 uploadFile 路由
func Use(api fiber.Router) {
	router := api.Group("/upload")
	// 上传文件
	router.Post("/file/*", mw.PermissProtectd, uploadFile)
	// 上传多个文件
	router.Post("/files/*", mw.PermissProtectd, uploadFiles)
	// 上传文件夹（支持chrome）
	router.Post("/folder", mw.PermissProtectd, uploadFolder)
}

// uploadFile 上传文件
func uploadFile(c *fiber.Ctx) error {
	var baseDir = c.Params("*", "data")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	return c.SaveFile(file, path.Join(global.Config.FirstStaticDir(), baseDir, file.Filename))
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
		fs.CreateDir(dir)
		err := c.SaveFile(file, dst)
		if err != nil {
			return err
		}
	}
	global.DB.Exec("select * from pipal.update_app_lib_version()")

	return res.JSON(c, true)
}

// uploadFolder 上传文件夹
func uploadFolder(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	for _, file := range files {
		fileName := file.Filename
		if strings.HasPrefix(fileName, "dist") {
			fileName = strings.Replace(fileName, "dist/", "", 1)
			dst := path.Join(global.Config.FirstStaticDir(), fileName)
			fs.CreateDir(dst)
			err := c.SaveFile(file, dst)
			if err != nil {
				return res.Error(c, "上传失败", err)
			}
		} else {
			return res.Error(c, "上传的文件夹不符合规范", nil)
		}
	}
	global.DB.Exec("select * from pipal.update_app_lib_version()")
	return res.JSON(c, true)
}
