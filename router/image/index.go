package image

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/pluto/ts/img"
)

// Use 重定向
func Use(api fiber.Router) {
	router := api.Group("/image")
	router.Get("/*", imageHandle)
}

// 处理图片
func imageHandle(c *fiber.Ctx) error {
	args := &ImageArgs{}
	err := c.QueryParser(args)
	if err != nil {
		return err
	}
	path := c.Params("*")
	//
	p, err := img.FromFile(filepath.Join("public", path))
	if err != nil {
		return err
	}
	// resize 缩放
	if args.W > 0 || args.H > 0 {
		p = p.Resize(args.W, args.H)
	}

	return res.Image(c, p)
}
