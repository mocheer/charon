package dmap

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/models/orm"
	"github.com/mocheer/charon/src/models/types"
	"github.com/mocheer/charon/src/router/store"
)

// Use 初始化 dmap 路由
func Use(api fiber.Router) {
	router := api.Group("/dmap")
	//
	router.Get("/image/:id/:z/:y:x", store.GlobalCache, imageHandle)
	router.Get("/layer/:id", store.GlobalCache, layerHandle)
	router.Get("/feature/:id", store.GlobalCache, featureHandle)
	router.Get("/identify/:id", store.GlobalCache, identifyHandle)

}

// imageHandle
func imageHandle(c *fiber.Ctx) error {
	//
	idParam := c.Params("id")
	zParam := c.Params("z")
	yParam := c.Params("y")
	xParam := c.Params("x")

	id, _ := strconv.Atoi(idParam)
	x, _ := strconv.Atoi(xParam)
	y, _ := strconv.Atoi(yParam)
	z, _ := strconv.Atoi(zParam)

	layer := NewDynamicLayer(id, &types.Tile{
		Z: z, Y: y, X: x,
	})

	if layer != nil {

	}

	return res.ResultError(c, "获取数据错误", nil)
}

func layerHandle(c *fiber.Ctx) error {
	id := c.Params("id")
	builder := &orm.SelectBuilder{}
	builder.Name = "layer"
	builder.Mode = "first"
	builder.Where = fmt.Sprintf("id=%s", id)
	builder.Select = "crs,extent,id,name,options,type,properties,array_to_json(array(select row_to_json(e) from (select * from pipal.dmap_layer where parent_id = 1)e)) as items"
	result := builder.Query()

	return res.ResultOK(c, result)
}

func featureHandle(c *fiber.Ctx) error {
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

	builder.Select = "layer_id,id,ST_AsGeoJson(geometry,4) as geometry,options,properties"
	result := builder.Query()
	return res.ResultOK(c, result)
}

func identifyHandle(c *fiber.Ctx) error {
	return nil
}
