package orm

import (
	"encoding/json"
	"strings"

	"github.com/mocheer/charon/src/db"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/models"
	"gorm.io/gorm"
)

// SelectBuilder 构建查询语句的参数
type SelectBuilder struct {
	Name   string `query:"name"`   // 表名
	Mode   string `query:"mode"`   // 数据返回模式 take | first | last | find | default
	Where  string `query:"where"`  // 查询 a=1 || {a:1}  => where a=1
	Not    string `query:"not"`    // 查询 a=1 || {a:1}  => where not a=1
	Select string `query:"select"` // 字段 a as b, sum(a) as b
	Limit  int    `query:"limit"`  // 限制数据量
	Offset int    `query:"offset"` // 偏移位置
	Order  string `query:"order"`  // 排序
	// 私有属性
	tx     *gorm.DB
	entity models.IEntity
}

// Model 设置模型
func (builder *SelectBuilder) Model() *gorm.DB {
	if builder.tx == nil && builder.Name != "" {
		builder.entity = models.NewTableStruct(builder.Name)
		builder.tx = global.Db.Model(builder.entity)
	}
	return builder.tx
}

// Query 执行查询
func (builder *SelectBuilder) Query() interface{} {
	query := builder.Model()
	//
	if builder.Select != "" {
		// 不止是字段选择，还是字段重命名，且支持函数调用
		query.Select(strings.Split(builder.Select, ","))
	}
	//
	if builder.Where != "" {
		var whereMap map[string]interface{}
		err := json.Unmarshal([]byte(builder.Where), &whereMap)
		if err == nil {
			query.Where(whereMap)
		} else {
			query.Where(builder.Where)
		}
	}
	if builder.Not != "" {
		var notMap map[string]interface{}
		err := json.Unmarshal([]byte(builder.Not), &notMap)
		if err == nil {
			query.Not(notMap)
		} else {
			query.Not(builder.Not)
		}
	}
	//
	if builder.Order != "" {
		query.Order(builder.Order)
	}

	if builder.Limit != 0 {
		query.Limit(builder.Limit)
	}
	//
	if builder.Offset != 0 {
		query.Offset(builder.Offset)
	}

	//
	switch builder.Mode {
	case "first":
		query.First(builder.entity)
		return builder.entity
	case "last":
		query.Last(&builder.entity)
		return builder.entity
	case "find":
		var result = models.NewTableStructArray(builder.Name)
		query.Find(result)
		return result
	default:
		var result []map[string]interface{}
		db.ScanIntoMap(query, &result)
		return result
	}
}
