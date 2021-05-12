package query

import (
	"encoding/json"

	"github.com/mocheer/charon/model"
	"github.com/mocheer/charon/mw"

	"github.com/mocheer/charon/model/orm"
	"github.com/mocheer/charon/model/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/mw/res"
)

// Use 初始化 query 路由
// @see http://apijson.cn/
func Use(api fiber.Router) {
	//
	router := api.Group("/query")
	// select
	router.Get("/:name", mw.Cache, Select)
	// insert 需要添加认证
	router.Put("/:name", mw.PermissProtectd, Insert)
	// update 需要添加认证
	router.Post("/:name", mw.PermissProtectd, Update)
	// delete 需要添加认证
	router.Delete("/:name", mw.PermissProtectd, Delete)
	// raw
	router.Get("/raw/:name", Raw)
	router.Post("/raw/:name", Raw)
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
	var entity = model.NewTableStruct(nameParam)
	c.BodyParser(entity)
	var db = global.DB.Model(entity)
	db.Create(entity)
	return res.JSON(c, true)
}

// Update 修改
func Update(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	var entity = model.NewTableStruct(nameParam)
	var db = global.DB.Model(entity)
	json.Unmarshal(c.Body(), entity)
	if whereQuery != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(whereQuery), &whereMap)
		if err == nil {
			db.Where(whereMap)
		} else {
			db.Where(whereQuery)
		}
	}
	db = db.Updates(entity)
	if db.RowsAffected > 0 {
		return res.JSON(c, true)
	} else {
		return res.Result(c, 500, false, "修改失败")
	}
}

// Delete
func Delete(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	//
	var entity = model.NewTableStruct(nameParam)
	json.Unmarshal(c.Body(), entity)
	var db = global.DB.Table(entity.TableName())
	//
	if whereQuery != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(whereQuery), &whereMap)
		if err == nil {
			db.Where(whereMap)
		} else {
			db.Where(whereQuery)
		}
	}
	db.Delete(entity)
	return res.JSON(c, true)
}

// Raw
func Raw(c *fiber.Ctx) error {
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
		db := global.DB.Raw(rawSQL.Text, namedArg)
		//
		global.DB.ScanIntoMap(db, &result)
	}

	return res.JSON(c, result)
}
