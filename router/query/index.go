package query

import (
	"encoding/json"
	"fmt"

	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/orm"
	"github.com/mocheer/charon/req"

	"github.com/mocheer/charon/orm/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/res"
)

// Use 初始化 query 路由
// @see http://apijson.cn/
func Use(api fiber.Router) {
	//
	router := api.Group("/query")
	// insert 需要添加认证
	router.Put("/:name", mw.PermissProtectd, Insert)
	// update 需要添加认证
	router.Post("/:name", mw.PermissProtectd, Update)
	// delete 需要添加认证
	router.Delete("/:name", mw.PermissProtectd, Delete)
	// new
	router.Get("/", mw.Cache, Select)
	// insert 需要添加认证
	router.Put("/", mw.PermissProtectd, Insert)
	// update 需要添加认证
	router.Post("/", mw.PermissProtectd, Update)
	// delete 需要添加认证
	router.Delete("/", mw.PermissProtectd, Delete)
	// raw
	router.Get("/raw/:name", Raw)
	router.Post("/raw/:name", Raw)
}

// Select
func Select(c *fiber.Ctx) error {
	var args = &orm.SelectArgs{}
	if err := c.QueryParser(args); err != nil {
		panic(fmt.Sprintf("参数有误：%s", err.Error()))
	}
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
	// table
	req.Engine().Model(nameParam).Where(whereQuery).Delete(c.Body()).Success()
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
