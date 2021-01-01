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

// Use 初始化 table 路由
func Use(api fiber.Router) {
	router := api.Group("/query")
	// s=>select
	router.Get("/s/:name", store.GlobalCache, selectTable)
	// i=>insert 需要添加认证
	router.Post("/i/:name", auth.GlobalProtected, insertTable)
	// u=>update 需要添加认证
	router.Post("/u/:name", auth.GlobalProtected, updateTable)
	// d=>delete 需要添加认证
	router.Post("/d/:name", auth.GlobalProtected, deleteTable)
	// raw
	router.Get("/raw/:name", queryRaw)
	router.Post("/raw/:name", queryRaw)

}

// matched, _ := regexp.MatchString(`pg_*`, nameParam)
// selectTable
func selectTable(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var modeQuery = c.Query("mode")   // take | first | last | find(default)
	var whereQuery = c.Query("where") // a=1 || {a:1}  => where a=1
	var notQuery = c.Query("not")     // a=1 || {a:1}  => where not a=1
	var selectQuery = c.Query("select")
	var limitQuery = c.Query("limit")
	var orderQuery = c.Query("order")
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
	case "last":
		query.Last(entity)
		return res.ResultOK(c, entity)
	default:
		db.ScanIntoMap(query, &result)
	}

	//

	return res.ResultOK(c, result)
}

func insertTable(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var entity = models.TableMapGenerate[nameParam]()
	json.Unmarshal(c.Body(), entity)
	var query = global.Db.Table(entity.TableName())
	query.Create(entity)
	return res.ResultOK(c, true)
}

func updateTable(c *fiber.Ctx) error {
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
	return res.ResultOK(c, true)
}

func deleteTable(c *fiber.Ctx) error {
	var nameParam = c.Params("name")
	var body = models.TableMapGenerate[nameParam]()
	json.Unmarshal(c.Body(), &body)
	global.Db.Delete(&body)
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
