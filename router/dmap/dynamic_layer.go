package dmap

import (
	"fmt"
	"runtime/debug"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/orm"
	"github.com/mocheer/charon/orm/tables"
	"github.com/mocheer/charon/req"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/js"
	"github.com/mocheer/pluto/ts/geois"
)

// createImageHandle
func createImageHandle(c *fiber.Ctx) error {
	//
	idParam := c.Params("id")
	zParam := c.Params("z")
	//
	id, _ := strconv.Atoi(idParam)
	z, _ := strconv.Atoi(zParam)
	//

	args := &orm.SelectArgs{}
	args.Name = "layer"
	args.Mode = "first"
	args.Where = fmt.Sprintf("id=%s", idParam)
	layerInfo := req.Engine().Query(args)
	fmt.Println(layerInfo)
	layer := layerInfo.(*tables.DmapLayer)
	//
	dynamicLayer := NewDynamicLayer(id, &geois.Tile{
		Z: z,
	})
	dynamicLayer.SetOptions(layer.Options)
	//
	if dynamicLayer != nil {
		args := &orm.SelectArgs{}
		args.Name = "feature"
		args.Mode = "find"
		args.Where = fmt.Sprintf("layer_id=%s", idParam)

		result := req.Engine().Query(args)
		features := result.(*[]tables.DmapFeature)
		for _, feature := range *features {
			dynamicLayer.Add(feature.Geometry)
		}
		dynamicLayer.Draw()
		cancelInterval := js.SetInterval(debug.FreeOSMemory, 60)
		defer cancelInterval()
		dynamicLayer.SaveTiles().Wait()
		if dynamicLayer.NumTile < 32 {
			data := dynamicLayer.GetData()
			return res.PNG(c, data)
		}
		debug.FreeOSMemory()
		return res.JSON(c, true)
	}

	return res.Error(c, "获取数据错误", nil)
}

// imageHandle
func imageTileHandle(c *fiber.Ctx) error {
	//
	idParam := c.Params("id")
	zParam := c.Params("z")
	yParam := c.Params("y")
	xParam := c.Params("x")
	//
	path := fmt.Sprintf(ImageTilePathFormat, idParam, zParam, yParam, xParam)
	if fs.IsExist(path) {
		return c.SendFile(path)
	}
	//
	id, _ := strconv.Atoi(idParam)
	x, _ := strconv.Atoi(xParam)
	y, _ := strconv.Atoi(yParam)
	z, _ := strconv.Atoi(zParam)
	//
	args := &orm.SelectArgs{}
	args.Name = "layer"
	args.Mode = "first"
	args.Where = fmt.Sprintf("id=%s", idParam)
	layerInfo := req.Engine().Query(args)
	layer := layerInfo.(*tables.DmapLayer)
	//
	dynamicLayer := NewDynamicLayer(id, &geois.Tile{
		Z: z, Y: y, X: x,
	})
	dynamicLayer.SetOptions(layer.Options)
	//
	if dynamicLayer != nil {
		args := &orm.SelectArgs{}
		args.Name = "feature"
		args.Mode = "find"
		args.Where = fmt.Sprintf("layer_id=%s", idParam)

		result := req.Engine().Query(args)
		features := result.(*[]tables.DmapFeature)

		for _, feature := range *features {
			dynamicLayer.Add(feature.Geometry)
		}
		dynamicLayer.Draw()
		fp := dynamicLayer.GetTile(x, y)
		// GC 垃圾回收太慢，需要手动释放
		debug.FreeOSMemory()
		return c.SendFile(fp)
	}

	return res.Error(c, "获取数据错误", nil)
}
