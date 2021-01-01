package db

import (
	"encoding/json"

	"github.com/mocheer/charon/src/core/fn"
	"github.com/mocheer/charon/src/global"
	"gorm.io/gorm"
)

// ScanIntoMap 扫描数据
func ScanIntoMap(query *gorm.DB, data *[]map[string]interface{}) {
	rows, err := query.Rows()
	result := *data
	if err == nil {
		for rows.Next() {
			rowData := map[string]interface{}{}
			err = global.Db.ScanRows(rows, &rowData)
			//
			if err == nil {
				columnTypes, _ := rows.ColumnTypes()
				for _, columnType := range columnTypes {
					// gorm-postgresql 底层的 DataType 并没有识别JSON，然后转成[]byte，而是用的string
					if columnType.DatabaseTypeName() == "JSON" {
						name := columnType.Name()
						if fn.IsString(rowData[name]) {
							var data interface{}
							err := json.Unmarshal([]byte(rowData[name].(string)), &data)
							if err == nil {
								rowData[name] = data
							}
						}
					}
				}
			}
			//
			result = append(result, rowData)
		}
	}
	*data = result
}
