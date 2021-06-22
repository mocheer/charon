package dmap

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/req"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/pluto/ts"
)

// layerHandle
func layerHandle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	result := req.Engine().
		Model("layer").
		Where(fmt.Sprintf("id=%s", idParam)).
		SelectTableAsJSON(ts.Map{"id": idParam}).
		First()
	//
	return res.JSON(c, result)
}
