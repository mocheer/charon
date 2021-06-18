package dmap

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/orm"
	"github.com/mocheer/charon/orm/model"
	"github.com/mocheer/charon/orm/tables"
	"github.com/mocheer/charon/req"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/js"
	"github.com/mocheer/pluto/ts"
	"github.com/mocheer/pluto/ts/geois"
)

// Use 初始化 dmap 路由
func Use(api fiber.Router) {
	router := api.Group("/dmap")
	router.Get("/image/:id/:z", mw.NewLimiter(limiter.Config{Max: 1, Expiration: 30 * time.Second}), createImageHandle)
	router.Get("/image/:id/:z/:y/:x", imageTileHandle)
	router.Get("/layer/:id", mw.Cache, layerHandle)
	router.Get("/feature/:id", mw.Cache, featureHandle)
	router.Get("/feature2/:id", mw.Cache, featureHandle2)
	router.Get("/identify/:id", mw.Cache, identifyHandle)
}

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

// featureHandle 要素服务
func featureHandle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	//
	result := &[]model.GeoFeature{}
	global.DB.Raw(`select row.geojson->>'type' as type , row.geojson->'coordinates' as coordinates , row.properties from (select st_asgeojson(geometry,4)::jsonb as geojson,properties from pipal.dmap_feature where layer_id = ?)row `, idParam).Scan(result)
	//
	return res.JSON(c, map[string]interface{}{
		"type":       "GeometryCollection",
		"geometries": result,
		"properties": struct{}{},
	})
}

// featureHandle2 要素服务
func featureHandle2(c *fiber.Ctx) error {
	id := c.Params("id")
	args := &orm.SelectArgs{}
	err := c.QueryParser(args)
	if err != nil {
		return err
	}
	args.Name = "feature"
	args.Mode = "find"

	if args.Where != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(args.Where), &whereMap)
		if err == nil {
			whereArray := []string{}
			for name, val := range whereMap {
				str := fmt.Sprintf("properties->>'%s'='%s'", name, val)
				whereArray = append(whereArray, str)
			}
			args.Where = fmt.Sprintf("layer_id=%s and %s", id, strings.Join(whereArray, " and "))
		} else {
			args.Where = fmt.Sprintf("layer_id=%s and %s", id, args.Where)
		}
	} else {
		args.Where = fmt.Sprintf("layer_id=%s", id)
	}
	//
	result := req.Engine().Query(args)
	return res.JSON(c, result)
}

// identifyHandle
func identifyHandle(c *fiber.Ctx) error {
	return nil
}
