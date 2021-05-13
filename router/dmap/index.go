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
	"github.com/mocheer/charon/model/orm"
	"github.com/mocheer/charon/model/tables"
	"github.com/mocheer/charon/model/types"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/js/window"
	"github.com/mocheer/pluto/ts"
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
	builder := &orm.SelectBuilder{}
	builder.Name = "layer"
	builder.Mode = "first"
	builder.Where = fmt.Sprintf("id=%s", idParam)
	layerInfo := builder.Query()
	layer := layerInfo.(*tables.DmapLayer)
	//
	dynamicLayer := NewDynamicLayer(id, &ts.Tile{
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
		cancelInterval := window.SetInterval(debug.FreeOSMemory, 60)
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
	builder := &orm.SelectBuilder{}
	builder.Name = "layer"
	builder.Mode = "first"
	builder.Where = fmt.Sprintf("id=%s", idParam)
	layerInfo := builder.Query()
	layer := layerInfo.(*tables.DmapLayer)
	//
	dynamicLayer := NewDynamicLayer(id, &ts.Tile{
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

	return res.Error(c, "获取数据错误", nil)
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
	return res.JSON(c, result)
}

// featureHandle 要素服务
func featureHandle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	result := &[]types.GeoFeature{}
	global.DB.Raw(`select row.geojson->>'type' as type , row.geojson->'coordinates' as coordinates , row.properties from (select st_asgeojson(geometry,4)::jsonb as geojson,properties from pipal.dmap_feature where layer_id = ?)row `, idParam).Scan(result)
	//
	return res.JSON(c, &map[string]interface{}{
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
	return res.JSON(c, result)
}

// identifyHandle
func identifyHandle(c *fiber.Ctx) error {
	return nil
}
