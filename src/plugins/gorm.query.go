package plugins

import (
	"encoding/json"

	"github.com/mocheer/pluto/fn"
	"gorm.io/gorm"
)

// UseGormQuery 弃用
func UseGormQuery(db *gorm.DB) {
	db.Callback().Query().Register("charon:query", charonQuery)
}

func charonQuery(db *gorm.DB) {
	if db.Statement.Schema == nil {
		value := db.Statement.ReflectValue.Interface().(map[string]interface{})
		for key, val := range value {
			// fmt.Println(key)
			if val != nil && fn.IsString(val) { //json字段的反射值仍然是字符串，实际上应该是[]byte更为合理
				var data interface{}
				err := json.Unmarshal([]byte(val.(string)), &data)
				if err == nil {
					value[key] = data
				}
			}
		}
	}
}
