package image

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/pluto/ts/img"
)

// Use 重定向
func Use(api fiber.Router) {
	router := api.Group("/image")
	router.Get("/*", getImage)
}

func getImage(c *fiber.Ctx) error {
	args := &GetImageArgs{}
	err := c.QueryParser(args)
	if err != nil {
		panic(err)
	}
	path := c.Params("*")

	if err != nil {
		return err
	}

	img.FromFile(path)

	return nil
}
