package dmap

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/models/orm"
	"github.com/mocheer/charon/src/models/tables"
	"github.com/mocheer/charon/src/models/types"
	"github.com/mocheer/charon/src/router/store"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/tm"
)

// Use 初始化 dmap 路由
func Use(api fiber.Router) {
	router := api.Group("/dmap")
	router.Get("/image/:id/:z", imageHandle)
	router.Get("/image/:id/:z/:y/:x", imageTileHandle)
	router.Get("/layer/:id", store.GlobalCache, layerHandle)
	router.Get("/feature/:id", store.GlobalCache, featureHandle)
	router.Get("/feature2/:id", store.GlobalCache, featureHandle2)
	router.Get("/identify/:id", store.GlobalCache, identifyHandle)
}

// imageHandle
func imageHandle(c *fiber.Ctx) error {
	//
	idParam := c.Params("id")
	zParam := c.Params("z")
	//
	id, _ := strconv.Atoi(idParam)
	z, _ := strconv.Atoi(zParam)
	//
	builder := &orm.SelectBuilder{}
	builder.Name = "layer"
	builder.Mode = "first"
	builder.Where = fmt.Sprintf("id=%s", idParam)
	layerInfo := builder.Query()
	layer := layerInfo.(*tables.DmapLayer)
	//
	dynamicLayer := NewDynamicLayer(id, &types.Tile{
		Z: z,
	})
	dynamicLayer.SetOptions(layer.Options)
	//
	if dynamicLayer != nil {
		builder := &orm.SelectBuilder{}
		builder.Name = "feature"
		builder.Mode = "find"
		builder.Where = fmt.Sprintf("layer_id=%s", idParam)

		result := builder.Query()
		features := result.(*[]tables.DmapFeature)

		for _, feature := range *features {
			dynamicLayer.Add(feature.Geometry)
		}
		dynamicLayer.Draw()
		cancelInterval := tm.SetInterval(debug.FreeOSMemory, 60)
		defer cancelInterval()
		dynamicLayer.SaveTiles().Wait()
		if dynamicLayer.NumTile < 32 {
			data := dynamicLayer.GetData()
			return res.ResultPNG(c, data)
		}
		debug.FreeOSMemory()
		return res.ResultOK(c, true)
	}

	return res.ResultError(c, "获取数据错误", nil)
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
	builder := &orm.SelectBuilder{}
	builder.Name = "layer"
	builder.Mode = "first"
	builder.Where = fmt.Sprintf("id=%s", idParam)
	layerInfo := builder.Query()
	layer := layerInfo.(*tables.DmapLayer)
	//
	dynamicLayer := NewDynamicLayer(id, &types.Tile{
		Z: z, Y: y, X: x,
	})
	dynamicLayer.SetOptions(layer.Options)
	//
	if dynamicLayer != nil {
		builder := &orm.SelectBuilder{}
		builder.Name = "feature"
		builder.Mode = "find"
		builder.Where = fmt.Sprintf("layer_id=%s", idParam)

		result := builder.Query()
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

	return res.ResultError(c, "获取数据错误", nil)
}

// layerHandle
func layerHandle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	builder := &orm.SelectBuilder{}
	builder.Name = "layer"
	builder.Mode = "first"
	builder.Where = fmt.Sprintf("id=%s", idParam)
	builder.Select = "crs,extent,id,name,options,type,properties,array_to_json(array(select row_to_json(e) from (select * from pipal.dmap_layer where parent_id = 1)e)) as items"
	result := builder.Query()
	//
	return res.ResultOK(c, result)
}

// featureHandle 要素服务
func featureHandle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	result := &[]types.GeoFeature{}
	global.Db.Raw(`select row.geojson->>'type' as type , row.geojson->'coordinates' as coordinates , row.properties from (select st_asgeojson(geometry,4)::jsonb as geojson,properties from pipal.dmap_feature where layer_id = ?)row `, idParam).Scan(result)
	//
	return res.ResultOK(c, &map[string]interface{}{
		"type":       "GeometryCollection",
		"geometries": result,
		"properties": struct{}{},
	})
}

// featureHandle2 要素服务
func featureHandle2(c *fiber.Ctx) error {
	id := c.Params("id")
	builder := &orm.SelectBuilder{}
	err := c.QueryParser(builder)
	if err != nil {
		return err
	}
	builder.Name = "feature"
	builder.Mode = "find"

	if builder.Where != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(builder.Where), &whereMap)
		if err == nil {
			whereArray := []string{}
			for name, val := range whereMap {
				str := fmt.Sprintf("properties->>'%s'='%s'", name, val)
				whereArray = append(whereArray, str)
			}
			builder.Where = fmt.Sprintf("layer_id=%s and %s", id, strings.Join(whereArray, " and "))
		} else {
			builder.Where = fmt.Sprintf("layer_id=%s and %s", id, builder.Where)
		}
	} else {
		builder.Where = fmt.Sprintf("layer_id=%s", id)
	}
	//
	result := builder.Query()
	return res.ResultOK(c, result)
}

// identifyHandle
func identifyHandle(c *fiber.Ctx) error {
	return nil
}
