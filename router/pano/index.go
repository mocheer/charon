package pano

import "github.com/gofiber/fiber/v2"

func Use(api fiber.Router) {
	router := api.Group("/pano")
	router.Post("/make", makePano)

}

// 制作全景图
func makePano(c *fiber.Ctx) error {
	return nil
}
