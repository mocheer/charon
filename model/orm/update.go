package orm

import (
	"encoding/json"

	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/model"
)

// UpdateBuilder 构建查询语句的参数
type UpdateBuilder struct {
	Name  string `query:"name"`  // 表名
	Where string `query:"where"` // 查询 a=1 || {a:1}  => where a=1
}

// Query 执行查询
func (builder *UpdateBuilder) Query() interface{} {
	var entity = model.NewTableStruct(builder.Name)
	//
	query := global.DB.Model(entity)
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
	return nil
}
