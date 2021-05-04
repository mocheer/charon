package model

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/models"
)

// Use 初始化 pipal 路由
func Use(api fiber.Router) {
	router := api.Group("/model")
	// query
	router.Get("/table/:name", queryTable)

}

func queryTable(c *fiber.Ctx) error {
	var result []map[string]interface{}
	var query = global.DB.Raw(`
SELECT 
tb.tablename as tablename,
a.attname AS columnname,
t.typname AS type
FROM
pg_class as c,
pg_attribute as a, 
pg_type as t,
(select tablename from pg_tables where schemaname = @schema) as tb
WHERE  a.attnum > 0 
and a.attrelid = c.oid
and a.atttypid = t.oid
and c.relname = tb.tablename 
order by tablename
`, map[string]interface{}{"schema": "pipal"})
	global.DB.ScanIntoMap(query, &result)
	var str = ""
	for tableName, ctMap := range result {
		str += fmt.Sprintf("type %d struct {\n", tableName)
		for colName, coltype := range ctMap {
			//translate to go struct foramt
			retype := models.PgTypeMap[coltype.(string)]
			if retype == "" {
				retype = coltype.(string)
			}
			jsonName := strings.ToLower(colName)
			vname := strings.ToUpper(jsonName[0:1]) + jsonName[1:]
			str += fmt.Sprintf("  %-10s %-15s `json:\"%s\"`\n", vname, retype, jsonName)
			//translate to typescript format
			retype = models.PgTypeMap[coltype.(string)]
			if retype == "" {
				retype = coltype.(string)
			}
		}
		str += "}\n"
	}
	return c.Send([]byte(str))
}
