package query

import (
	"encoding/json"
	"strconv"
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
	var nameParam = c.Params("name")    // 表名
	var modeQuery = c.Query("mode")     // 数据返回模式 take | first | last | find(default)
	var whereQuery = c.Query("where")   // 查询 a=1 || {a:1}  => where a=1
	var notQuery = c.Query("not")       // 查询 a=1 || {a:1}  => where not a=1
	var selectQuery = c.Query("select") // 字段
	var limitQuery = c.Query("limit")   //
	// var offsetQuery = c.Query("offset") // 偏移
	var orderQuery = c.Query("order") // 排序
	var result []map[string]interface{}
	var query = global.Db
	var entity = models.TableMapGenerate[nameParam]()
	//
	query = query.Table(entity.TableName())
	//
	if selectQuery != "" {
		// 不止是字段选择，还是字段重命名，且支持函数调用
		query.Select(strings.Split(selectQuery, ","))

	}

	if whereQuery != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(whereQuery), &whereMap)
		if err == nil {
			query.Where(whereMap)
		} else {
			query.Where(whereQuery)
		}
	}
	if notQuery != "" {
		var notMap map[string]interface{}
		err := json.Unmarshal([]byte(notQuery), &notMap)
		if err == nil {
			query.Not(notMap)
		} else {
			query.Not(notQuery)
		}
	}
	//
	if orderQuery != "" {
		query.Order(orderQuery)
	}
	//
	if limitQuery != "" {
		limit, _ := strconv.Atoi(limitQuery)
		query.Limit(limit)
	}

	//
	switch modeQuery {
	case "first":
		query.First(entity)
		return res.ResultOK(c, entity)
		// return res.ResultOK(c, nil)
	case "last":
		query.Last(entity)
		return res.ResultOK(c, entity)
	case "find":
		query.Find(result)
		break
	default:
		db.ScanIntoMap(query, &result)
	}

	return res.ResultOK(c, result)
}

func queryInsert(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var entity = models.TableMapGenerate[nameParam]()
	json.Unmarshal(c.Body(), entity)
	var query = global.Db.Table(entity.TableName())
	query.Create(entity)
	return res.ResultOK(c, true)
}

func queryUpdate(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")

	var body map[string]interface{}
	var entity = models.TableMapGenerate[nameParam]()
	var query = global.Db.Table(entity.TableName())
	json.Unmarshal(c.Body(), &body)
	if whereQuery != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(whereQuery), &whereMap)
		if err == nil {
			query.Where(whereMap)
		} else {
			query.Where(whereQuery)
		}
	}

	query.Updates(body)
	return res.ResultOK(c, query.RowsAffected > 1)
}

// queryDelete
func queryDelete(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var whereQuery = c.Query("where")
	//
	var entity = models.TableMapGenerate[nameParam]()
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
