package upload

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Use 初始化 uploadFile 路由
func Use(api fiber.Router) {
	router := api.Group("/upload")
	router.Post("/file", uploadFile)
	router.Post("/files", uploadFiles)
}

// uploadFile
func uploadFile(c *fiber.Ctx) error {
	// Get first file from form field "document":
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	// Save file to root directory:
	return c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))
}

// uploadFiles
func uploadFiles(c *fiber.Ctx) error {
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
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
		// => "tutorial.pdf" 360641 "application/pdf"

		// Save the files to disk:
		err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))

		// Check for errors
		if err != nil {
			return err
		}
	}
	return nil
}
