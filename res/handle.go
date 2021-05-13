package res

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/pluto/fs"
)

// HandleTextFile
func HandleTextFile(path string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := fs.MustReadText(path)
		return JSON(c, data)
	}
}
