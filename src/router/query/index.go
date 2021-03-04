package query

import (
	"encoding/json"
	"strings"

	"github.com/mocheer/charon/src/models"
	"github.com/mocheer/charon/src/router/auth"
	"github.com/mocheer/charon/src/router/store"

	"github.com/mocheer/charon/src/core/db"
	"github.com/mocheer/charon/src/models/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/global"
)

// Use 初始化 query 路由
// @see http://apijson.cn/
func Use(api fiber.Router) {
	//
	router := api.Group("/query")
	// select
	router.Get("/:name", store.GlobalCache, querySeclect)
	// insert 需要添加认证
	router.Put("/:name", auth.GlobalProtected, auth.PermissProtectd, queryInsert)
	// update 需要添加认证
	router.Post("/:name", auth.GlobalProtected, auth.PermissProtectd, queryUpdate)
	// delete 需要添加认证
	router.Delete("/:name", auth.GlobalProtected, auth.PermissProtectd, queryDelete)
	// raw
	router.Get("/raw/:name", queryRaw)
	router.Post("/raw/:name", queryRaw)
}

// matched, _ := regexp.MatchString(`pg_*`, nameParam)
// querySeclect
func querySeclect(c *fiber.Ctx) error {
	var nameParam = c.Params("name") // 表名
	var params = &SelectParams{}
	if err := c.QueryParser(params); err != nil {
		return err
	}

	var query = global.Db
	var entity = models.NewTableStruct(nameParam)
	//
	query = query.Model(entity)
	//
	if params.Select != "" {
		// 不止是字段选择，还是字段重命名，且支持函数调用
		query.Select(strings.Split(params.Select, ","))
	}
	//
	if params.Where != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(params.Where), &whereMap)
		if err == nil {
			query.Where(whereMap)
		} else {
			query.Where(params.Where)
		}
	}
	if params.Not != "" {
		var notMap map[string]interface{}
		err := json.Unmarshal([]byte(params.Not), &notMap)
		if err == nil {
			query.Not(notMap)
		} else {
			query.Not(params.Not)
		}
	}
	//
	if params.Order != "" {
		query.Order(params.Order)
	}

	if params.Limit != 0 {
		query.Limit(params.Limit)
	}
	//
	if params.Offset != 0 {
		query.Offset(params.Offset)
	}

	//
	switch params.Mode {
	case "first":
		query.First(entity)
		return res.ResultOK(c, entity)
	case "last":
		query.Last(&entity)
		return res.ResultOK(c, entity)
	case "find":
		var result = models.NewTableStructArray(nameParam)
		query.Find(result)
		return res.ResultOK(c, result)
	default:
		var result []map[string]interface{}
		db.ScanIntoMap(query, &result)
		return res.ResultOK(c, result)
	}
}

func queryInsert(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var entity = models.NewTableStruct(nameParam)
	c.BodyParser(entity)
	var query = global.Db.Model(entity)
	query.Create(entity)
	return res.ResultOK(c, true)
}

func queryUpdate(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	var entity = models.NewTableStruct(nameParam)
	var query = global.Db.Model(entity)
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

	query.Updates(entity)
	return res.ResultOK(c, query.RowsAffected > 1)
}

// queryDelete
func queryDelete(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	//
	var entity = models.NewTableStruct(nameParam)
	json.Unmarshal(c.Body(), entity)
	var query = global.Db.Table(entity.TableName())
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
	return res.ResultOK(c, true)
}

// queryRaw
func queryRaw(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var rawSQL tables.RawSQL
	var result []map[string]interface{}
	var err = global.Db.Where(&tables.RawSQL{Name: nameParam}).First(&rawSQL)
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
		query := global.Db.Raw(rawSQL.Text, namedArg)
		//
		db.ScanIntoMap(query, &result)
	}

	return res.ResultOK(c, result)
}
