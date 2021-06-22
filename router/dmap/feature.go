package dmap

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/orm"
	"github.com/mocheer/charon/orm/model"
	"github.com/mocheer/charon/req"
	"github.com/mocheer/charon/res"
)

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
