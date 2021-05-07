package global

import (
	"encoding/json"
	"time"

	"github.com/mocheer/pluto/fn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBase struct {
	*gorm.DB
}

// Open 连接数据库
func openDB() (database *DataBase, err error) {
	if Config.DSN == "" {
		return
	}
	db, err := gorm.Open(postgres.Open(Config.DSN), &gorm.Config{AllowGlobalUpdate: false})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(8)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(128)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	database = &DataBase{db}
	return
}

// ScanIntoMap 扫描数据
func (db *DataBase) ScanIntoMap(query *gorm.DB, data *[]map[string]interface{}) {
	rows, err := query.Rows()
	result := *data
	if err == nil {
		for rows.Next() {
			rowData := map[string]interface{}{}
			err = DB.ScanRows(rows, &rowData)
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
