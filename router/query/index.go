package query

import (
	"encoding/json"

	"github.com/mocheer/charon/model"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/req"

	"github.com/mocheer/charon/model/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/res"
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
	args := req.MustParseSelectArgs(c)
	args.Name = c.Params("name")
	result := req.Engine().Query(args)
	return res.JSON(c, result)
}

// Insert 增加
func Insert(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	successed := req.Engine().Model(nameParam).Create(c.Body()).Success()
	return res.JSON(c, successed)
}

// Update 修改
func Update(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	successed := req.Engine().Model(nameParam).Where(whereQuery).Updates(c.Body()).Success()
	return res.JSON(c, successed)
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

		req.Engine().Raw(rawSQL.Text, namedArg).ScanIntoMap(&result)
	}

	return res.JSON(c, result)
}
