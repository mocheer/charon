package mw

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// HideJSMap 不显示前端中的js.map资源
func HideJSMap(c *fiber.Ctx) error {
	if strings.HasSuffix(c.Path(), ".js.map") {
		return c.SendString(`{
		"version": 3,
		"sources": [],
		"names": [],
		"mappings": "",
		"file": "",
		"sourcesContent": [],
		"sourceRoot": ""
	}`)
	}
	return c.Next()
}
