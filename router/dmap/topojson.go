package dmap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/xena/gd"
	"github.com/rubenv/topojson"
)

func Geo2TopojsonHandle(ctx *fiber.Ctx) error {
	return res.JSON(ctx, Geo2Topojson(ctx.Body()))
}

func Geo2Topojson(in []byte) *topojson.Topology {
	topo := gd.NewTopology(in, &topojson.TopologyOptions{
		PreQuantize:  1000000,
		PostQuantize: 10000,
	})
	return topo
}
