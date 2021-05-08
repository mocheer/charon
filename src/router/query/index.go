package query

import (
	"encoding/json"

	"github.com/mocheer/charon/src/models"
	"github.com/mocheer/charon/src/mw"

	"github.com/mocheer/charon/src/models/orm"
	"github.com/mocheer/charon/src/models/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/mw/res"
)

// Use 初始化 query 路由
// @see http://apijson.cn/
func Use(api fiber.Router) {
	//
	router := api.Group("/query")
	// select
	router.Get("/:name", mw.Cache, Select)
	// insert 需要添加认证
	router.Put("/:name", mw.Protector, mw.PermissProtectd, Insert)
	// update 需要添加认证
	router.Post("/:name", mw.Protector, mw.PermissProtectd, Update)
	// delete 需要添加认证
	router.Delete("/:name", mw.Protector, mw.PermissProtectd, queryDelete)
	// raw
	router.Get("/raw/:name", queryRaw)
	router.Post("/raw/:name", queryRaw)
}

// matched, _ := regexp.MatchString(`pg_*`, nameParam)
// Select
func Select(c *fiber.Ctx) error {
	var builder = &orm.SelectBuilder{}
	if err := c.QueryParser(builder); err != nil {
		return err
	}
	builder.Name = c.Params("name")
	result := builder.Query()
	return res.JSON(c, result)
}

// Insert 增加
func Insert(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var entity = models.NewTableStruct(nameParam)
	c.BodyParser(entity)
	var query = global.DB.Model(entity)
	query.Create(entity)
	return res.JSON(c, true)
}

// Update 修改
func Update(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	var entity = models.NewTableStruct(nameParam)
	var query = global.DB.Model(entity)
	json.Unmarshal(c.Body(), entity)
	if whereQuery != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(whereQuery), &whereMap)
		if err == nil {
			query.Where(whereMap)
		} else {
			query.Where(whereQuery)
		}
	}

	query.Debug().Updates(entity)
	if query.RowsAffected > 1 {
		return res.JSON(c, true)
	} else {
		return res.Result(c, 500, false, "")
	}

}

// queryDelete
func queryDelete(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	//
	var entity = models.NewTableStruct(nameParam)
	json.Unmarshal(c.Body(), entity)
	var query = global.DB.Table(entity.TableName())
	//
	if whereQuery != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(whereQuery), &whereMap)
		if err == nil {
			query.Where(whereMap)
		} else {
			query.Where(whereQuery)
		}
	}
	query.Delete(entity)
	return res.JSON(c, true)
}

// queryRaw
func queryRaw(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var rawSQL tables.RawSQL
	var result []map[string]interface{}
	var err = global.DB.Where(&tables.RawSQL{Name: nameParam}).First(&rawSQL)
	if err != nil {
		var body map[string]interface{}
		var namedArg map[string]interface{}
		if c.Method() == "GET" {
			json.Unmarshal([]byte(c.Query("namedArg")), &namedArg)
		} else {
			// fetch POST contentType = applcation/json
			json.Unmarshal(c.Body(), &body)
			if body["namedArg"] != nil {
				namedArg = body["namedArg"].(map[string]interface{})
			}
		}
		query := global.DB.Raw(rawSQL.Text, namedArg)
		//
		global.DB.ScanIntoMap(query, &result)
	}

	return res.JSON(c, result)
}
