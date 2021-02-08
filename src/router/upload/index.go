package upload

import (
	"fmt"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/fn"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/router/auth"
)

// Use 初始化 uploadFile 路由
func Use(api fiber.Router) {
	router := api.Group("/upload")
	// 上传文件
	router.Post("/file", auth.GlobalProtected, auth.PermissProtectd, uploadFile)
	// 上传多个文件
	router.Post("/files", auth.GlobalProtected, auth.PermissProtectd, uploadFiles)
	// 上传文件夹（支持chrome）
	router.Post("/folder", auth.GlobalProtected, auth.PermissProtectd, uploadFolder)
}

// uploadFile 上传文件
func uploadFile(c *fiber.Ctx) error {
	// Get first file from form field "document":
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	// Save file to root directory:
	return c.SaveFile(file, fmt.Sprintf("./data/%s", file.Filename))
}

// uploadFiles 上传多个文件
func uploadFiles(c *fiber.Ctx) error {
	var baseDir = c.Query("*")
	if baseDir == "" {
		baseDir = "data"
	}
	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	// => *multipart.Form

	// Get all files from "documents" key:
	files := form.File["files"]
	// => []*multipart.FileHeader

	// Loop through files:
	for _, file := range files {
		// fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
		// => "tutorial.pdf" 360641 "application/pdf"

		// Save the files to disk:
		dst := path.Join(baseDir, file.Filename)
		dir := path.Dir(dst)
		fn.MkdirNotExist(dir)
		err := c.SaveFile(file, dst)

		// Check for errors
		if err != nil {
			return err
		}
	}
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
			dir := path.Dir(dst)
			fn.MkdirNotExist(dir)
			err := c.SaveFile(file, dst)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("error dir:", fileName)
		}
	}
	global.Db.Exec("select * from pipal.update_app_lib_version()")
	return res.ResultOK(c, true)
}
